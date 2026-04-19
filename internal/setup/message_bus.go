package setup

import (
	"fmt"
	"strings"
	"time"

	"codebase-app/internal/framework/secondary/publisher/nats"
	"codebase-app/internal/framework/secondary/publisher/postgres"
	"codebase-app/internal/infrastructure/config"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

func NewMessagePublisher(db *sqlx.DB) (integrationPorts.MessagePublisher, error) {
	switch normalizedMessageBusDriver() {
	case "postgres":
		return postgres.NewPublisher(db), nil
	case "nats":
		return nats.NewPublisher(config.Envs.EmailVerificationQueueNats.NatsURL)
	default:
		return nil, fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
	}
}

func NewMessageSubscriptionManager(db *sqlx.DB) (integrationPorts.MessageSubscriptionManager, error) {
	switch normalizedMessageBusDriver() {
	case "postgres":
		return postgres.NewSubscriptionManager(db, postgres.Config{
			PollInterval: time.Duration(config.Envs.MessageBus.PostgresPollIntervalMS) * time.Millisecond,
			BatchSize:    config.Envs.MessageBus.PostgresBatchSize,
			RetryDelay:   time.Duration(config.Envs.MessageBus.RetryDelaySeconds) * time.Second,
			MaxAttempts:  config.Envs.MessageBus.MaxAttempts,
		}), nil
	case "nats":
		return nats.NewSubscriptionManager(config.Envs.EmailVerificationQueueNats.NatsURL)
	default:
		return nil, fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
	}
}

func normalizedMessageBusDriver() string {
	driver := strings.TrimSpace(strings.ToLower(config.Envs.MessageBus.Driver))
	if driver == "" {
		return "postgres"
	}

	return driver
}
