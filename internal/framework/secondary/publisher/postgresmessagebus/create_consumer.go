package postgresmessagebus

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

type consumerManager struct {
	db  *sqlx.DB
	cfg Config
}

var _ integrationPorts.MessageConsumerManager = &consumerManager{}

func NewConsumerManager(db *sqlx.DB, cfg Config) integrationPorts.MessageConsumerManager {
	return &consumerManager{
		db:  db,
		cfg: cfg.normalized(),
	}
}

func (m *consumerManager) CreateConsumer(ctx context.Context, cfg integrationPorts.MessageBusConsumerConfig) (integrationPorts.MessageBusConsumer, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgresmessagebus:create_consumer:CreateConsumer")
	defer span.End()

	return &postgresConsumer{
		db:  m.db,
		cfg: m.cfg,
		def: cfg,
	}, nil
}

func (m *consumerManager) Close() error {
	return nil
}
