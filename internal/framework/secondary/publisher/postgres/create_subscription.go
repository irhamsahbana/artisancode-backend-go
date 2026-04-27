package postgres

import (
	"context"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/jmoiron/sqlx"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/integration"
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

	subscriber, err := newWatermillSubscriber(
		watermillSQL.BeginnerFromStdSQL(m.db.DB),
		m.cfg,
		cfg,
	)
	if err != nil {
		return nil, err
	}

	return &postgresConsumer{
		db:         m.db,
		cfg:        m.cfg,
		def:        cfg,
		subscriber: subscriber,
	}, nil
}

func (m *consumerManager) Close() error {
	return nil
}
