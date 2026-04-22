package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type postgresConsumer struct {
	db  *sqlx.DB
	cfg Config
	def integrationPorts.MessageBusSubscriptionConfig
}

var _ integrationPorts.MessageBusSubscription = &postgresConsumer{}

func (c *postgresConsumer) Consume(ctx context.Context, handler func(context.Context, integrationPorts.MessageBusMessage)) (integrationPorts.MessageBusSubscriptionContext, error) {
	ctx, cancel := context.WithCancel(ctx)

	consumeCtx := &postgresConsumeContext{
		cancel: cancel,
	}
	consumeCtx.wg.Add(1)

	go func() {
		defer consumeCtx.wg.Done()
		c.consumeLoop(ctx, handler)
	}()

	return consumeCtx, nil
}

func (c *postgresConsumer) consumeLoop(ctx context.Context, handler func(context.Context, integrationPorts.MessageBusMessage)) {
	ticker := time.NewTicker(c.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed, err := c.consumeBatch(ctx, handler)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to consume message queue batch")
		}

		if processed > 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (c *postgresConsumer) consumeBatch(ctx context.Context, handler func(context.Context, integrationPorts.MessageBusMessage)) (int, error) {
	messages := make([]*postgresMessage, 0, c.cfg.BatchSize)
	for range c.cfg.BatchSize {
		msg, err := c.claimMessage(ctx)
		if err != nil {
			return len(messages), err
		}
		if msg == nil {
			break
		}

		messages = append(messages, msg)
	}

	for _, msg := range messages {
		c.handleMessage(ctx, handler, msg)
	}

	return len(messages), nil
}

func (c *postgresConsumer) handleMessage(ctx context.Context, handler func(context.Context, integrationPorts.MessageBusMessage), msg *postgresMessage) {
	msgCtx := ctx
	cancel := func() {}
	if c.cfg.ProcessingTimeout > 0 {
		msgCtx, cancel = context.WithTimeout(ctx, c.cfg.ProcessingTimeout)
	}
	defer cancel()

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Ctx(msgCtx).Error().
				Any("panic", recovered).
				Bytes("stack", debug.Stack()).
				Str("subject", msg.Subject()).
				Msg("Message handler panicked")

			reason := fmt.Sprintf("handler panic: %v", recovered)
			if err := msg.Nak(msgCtx, reason); err != nil {
				log.Ctx(msgCtx).Error().Err(err).Str("subject", msg.Subject()).Msg("Failed to negative-acknowledge panicked message")
			}
		}
	}()

	handler(msgCtx, msg)
}

func (c *postgresConsumer) claimMessage(ctx context.Context) (*postgresMessage, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:consume:claimMessage")
	defer span.End()

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to begin message queue transaction")
		return nil, err
	}
	defer tx.Rollback()

	var row struct {
		ID          string `db:"id"`
		Subject     string `db:"subject"`
		PayloadJSON string `db:"payload_json"`
		HeadersJSON string `db:"headers_json"`
	}
	processingCutoff := time.Now().UTC().Add(-c.cfg.ProcessingTimeout)

	query := `
		SELECT id, subject, payload_json::text AS payload_json, headers_json::text AS headers_json
		FROM message_queue
		WHERE deleted_at IS NULL
		  AND (
			(status = 'pending' AND available_at <= CURRENT_TIMESTAMP)
			OR
			(status = 'processing' AND locked_at IS NOT NULL AND locked_at <= ?)
		  )
		  AND ` + buildSubjectPredicate(c.def.Subjects) + `
		ORDER BY COALESCE(available_at, locked_at) ASC, created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`

	err = tx.GetContext(ctx, &row, c.db.Rebind(query), processingCutoff)
	if err != nil {
		if err == sql.ErrNoRows {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Ctx(ctx).Error().Err(commitErr).Msg("Failed to commit empty queue claim")
				return nil, commitErr
			}

			return nil, nil
		}

		log.Ctx(ctx).Error().Err(err).Msg("Failed to claim queued message")
		return nil, err
	}

	updateQuery := `
		UPDATE message_queue
		SET status = 'processing',
			consumer_name = ?,
			locked_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, c.db.Rebind(updateQuery), c.def.ConsumerName, row.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to mark queued message as processing")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to commit queued message claim")
		return nil, err
	}

	headers := map[string][]string{}
	err = json.Unmarshal([]byte(row.HeadersJSON), &headers)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal queued message headers")
		return nil, err
	}

	return &postgresMessage{
		db:          c.db,
		id:          row.ID,
		subject:     row.Subject,
		data:        []byte(row.PayloadJSON),
		headers:     headers,
		retryDelay:  c.cfg.RetryDelay,
		maxAttempts: c.cfg.MaxAttempts,
	}, nil
}

type postgresConsumeContext struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var _ integrationPorts.MessageBusSubscriptionContext = &postgresConsumeContext{}

func (c *postgresConsumeContext) Stop() {
	c.cancel()
	c.wg.Wait()
}

func buildSubjectPredicate(subjects []string) string {
	if len(subjects) == 0 {
		return "1 = 1"
	}

	predicates := make([]string, 0, len(subjects))
	for _, subject := range subjects {
		if strings.HasSuffix(subject, ".>") {
			prefix := strings.TrimSuffix(subject, ">")
			predicates = append(predicates, fmt.Sprintf("subject LIKE '%s%%'", prefix))
			continue
		}

		predicates = append(predicates, fmt.Sprintf("subject = '%s'", subject))
	}

	return "(" + strings.Join(predicates, " OR ") + ")"
}
