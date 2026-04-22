package postgres

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type postgresMessage struct {
	db          *sqlx.DB
	id          string
	subject     string
	data        []byte
	headers     map[string][]string
	retryDelay  time.Duration
	maxAttempts int
}

var _ integrationPorts.MessageBusMessage = &postgresMessage{}

func (m *postgresMessage) Subject() string {
	return m.subject
}

func (m *postgresMessage) Data() []byte {
	return m.data
}

func (m *postgresMessage) Headers() map[string][]string {
	return m.headers
}

func (m *postgresMessage) Ack(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:message:Ack")
	defer span.End()

	query := `
		UPDATE message_queue
		SET status = 'processed',
			processed_at = CURRENT_TIMESTAMP,
			locked_at = NULL,
			last_error = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.db.ExecContext(ctx, m.db.Rebind(query), m.id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to acknowledge queued message")
		return err
	}

	return nil
}

func (m *postgresMessage) Nak(ctx context.Context, reason string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:message:Nak")
	defer span.End()

	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to begin negative-acknowledge transaction")
		return err
	}
	defer tx.Rollback()

	var row struct {
		Subject      string         `db:"subject"`
		PayloadJSON  string         `db:"payload_json"`
		HeadersJSON  string         `db:"headers_json"`
		ConsumerName sql.NullString `db:"consumer_name"`
		AttemptCount int            `db:"attempt_count"`
		CreatedAt    time.Time      `db:"created_at"`
	}

	selectQuery := `
		SELECT subject,
			payload_json::text AS payload_json,
			headers_json::text AS headers_json,
			consumer_name,
			attempt_count,
			created_at
		FROM message_queue
		WHERE id = ?
		FOR UPDATE
	`

	err = tx.GetContext(ctx, &row, m.db.Rebind(selectQuery), m.id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to load queued message before negative-acknowledge")
		return err
	}

	nextAttemptCount := row.AttemptCount + 1

	if m.maxAttempts > 0 && nextAttemptCount >= m.maxAttempts {
		insertDeadLetterQuery := `
			INSERT INTO message_queue_dead_letters (
				message_queue_id,
				subject,
				payload_json,
				headers_json,
				consumer_name,
				attempt_count,
				failure_reason,
				queued_at,
				dead_lettered_at,
				created_at
			) VALUES (?, ? , ?::jsonb, ?::jsonb, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`

		_, err = tx.ExecContext(
			ctx,
			m.db.Rebind(insertDeadLetterQuery),
			m.id,
			row.Subject,
			row.PayloadJSON,
			row.HeadersJSON,
			row.ConsumerName,
			nextAttemptCount,
			reason,
			row.CreatedAt,
		)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to insert queued message into dead letters")
			return err
		}

		deleteQuery := `DELETE FROM message_queue WHERE id = ?`
		_, err = tx.ExecContext(ctx, m.db.Rebind(deleteQuery), m.id)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to delete dead-lettered queued message")
			return err
		}
	} else {
		updateQuery := `
			UPDATE message_queue
			SET status = 'pending',
				locked_at = NULL,
				available_at = ?,
				attempt_count = ?,
				last_error = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`

		_, err = tx.ExecContext(
			ctx,
			m.db.Rebind(updateQuery),
			time.Now().UTC().Add(m.retryDelay),
			nextAttemptCount,
			reason,
			m.id,
		)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to negative-acknowledge queued message")
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to commit negative-acknowledge transaction")
		return err
	}

	return nil
}
