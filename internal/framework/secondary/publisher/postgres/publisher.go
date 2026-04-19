package postgres

import (
	"context"
	"encoding/json"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type publisher struct {
	db *sqlx.DB
}

var _ integrationPorts.MessagePublisher = &publisher{}

func NewPublisher(db *sqlx.DB) integrationPorts.MessagePublisher {
	return &publisher{db: db}
}

func (p *publisher) PublishJSON(ctx context.Context, subject string, payload any) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:publisher:PublishJSON")
	defer span.End()

	data, err := json.Marshal(payload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal message payload")
		return err
	}

	headers := map[string][]string{}
	otel.GetTextMapPropagator().Inject(ctx, mapHeaderCarrier(headers))

	headersJSON, err := json.Marshal(headers)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal message headers")
		return err
	}

	query := `
		INSERT INTO message_queue (
			subject,
			payload_json,
			headers_json,
			status,
			available_at,
			created_at,
			updated_at
		) VALUES (?, ?::jsonb, ?::jsonb, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = p.db.ExecContext(ctx, p.db.Rebind(query), subject, string(data), string(headersJSON))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to insert message into queue")
		return err
	}

	return nil
}

func (p *publisher) Close() error {
	return nil
}
