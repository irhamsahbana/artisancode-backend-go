package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"

	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/integration"
)

type publisher struct {
	db        *sqlx.DB
	publisher *watermillSQL.Publisher
}

var _ integrationPorts.MessagePublisher = &publisher{}

func NewPublisher(db *sqlx.DB) integrationPorts.MessagePublisher {
	watermillPublisher, err := newWatermillPublisher(watermillSQL.BeginnerFromStdSQL(db.DB))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Watermill SQL publisher")
	}

	return &publisher{
		db:        db,
		publisher: watermillPublisher,
	}
}

func (p *publisher) PublishJSON(ctx context.Context, subject string, payload any) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:publisher:PublishJSON")
	defer span.End()

	data, err := json.Marshal(payload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to marshal message payload")
		return err
	}

	metadata := message.Metadata{}
	metadata.Set("topic", subject)
	otel.GetTextMapPropagator().Inject(ctx, metadataHeaderCarrier(metadata))

	exec := postgresTx.ExecutorFromContext(ctx, p.db)
	watermillPublisher := p.publisher
	if tx, ok := exec.(*sqlx.Tx); ok {
		watermillPublisher, err = newWatermillPublisher(watermillSQL.TxFromStdSQL(tx.Tx))
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to initialize transactional Watermill SQL publisher")
			return err
		}
		defer watermillPublisher.Close()
	}

	msg := message.NewMessageWithContext(ctx, watermill.NewUUID(), data)
	msg.Metadata = metadata

	err = watermillPublisher.Publish(subject, msg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("topic", subject).Msg("Failed to publish message to Watermill SQL")
		return err
	}

	return nil
}

func (p *publisher) Close() error {
	if p.publisher == nil {
		return nil
	}

	return p.publisher.Close()
}

func newWatermillPublisher(db watermillSQL.ContextExecutor) (*watermillSQL.Publisher, error) {
	if db == nil {
		return nil, fmt.Errorf("watermill sql publisher db is nil")
	}

	return watermillSQL.NewPublisher(
		db,
		watermillSQL.PublisherConfig{
			SchemaAdapter:        watermillSQL.DefaultPostgreSQLSchema{},
			AutoInitializeSchema: false,
		},
		watermill.NopLogger{},
	)
}
