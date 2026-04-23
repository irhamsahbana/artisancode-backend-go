package postgres

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

type consumerManager struct {
	db  *sqlx.DB
	cfg Config
}

var _ integrationPorts.MessageSubscriptionManager = &consumerManager{}

func NewSubscriptionManager(db *sqlx.DB, cfg Config) integrationPorts.MessageSubscriptionManager {
	return &consumerManager{
		db:  db,
		cfg: cfg.normalized(),
	}
}

func (m *consumerManager) CreateSubscription(ctx context.Context, cfg common.MessageBusSubscriptionConfig) (integrationPorts.MessageBusSubscription, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:create_subscription:CreateSubscription")
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
