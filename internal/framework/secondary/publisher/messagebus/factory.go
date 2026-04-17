package messagebus

import (
	"fmt"
	"strings"
	"time"

	"codebase-app/internal/framework/secondary/publisher/natsjetstream"
	"codebase-app/internal/framework/secondary/publisher/postgresmessagebus"
	"codebase-app/internal/infrastructure/config"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/jmoiron/sqlx"
)

func NewPublisher(db *sqlx.DB) (integrationPorts.MessagePublisher, error) {
	switch normalizedDriver() {
	case "postgres":
		return postgresmessagebus.NewPublisher(db), nil
	case "nats":
		return natsjetstream.NewPublisher(config.Envs.EmailVerificationQueueNats.NatsURL)
	default:
		return nil, fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
	}
}

func NewConsumerManager(db *sqlx.DB) (integrationPorts.MessageConsumerManager, error) {
	switch normalizedDriver() {
	case "postgres":
		return postgresmessagebus.NewConsumerManager(db, postgresmessagebus.Config{
			PollInterval: time.Duration(config.Envs.MessageBus.PostgresPollIntervalMS) * time.Millisecond,
			BatchSize:    config.Envs.MessageBus.PostgresBatchSize,
			RetryDelay:   time.Duration(config.Envs.MessageBus.RetryDelaySeconds) * time.Second,
			MaxAttempts:  config.Envs.MessageBus.MaxAttempts,
		}), nil
	case "nats":
		return natsjetstream.NewConsumerManager(config.Envs.EmailVerificationQueueNats.NatsURL)
	default:
		return nil, fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
	}
}

func normalizedDriver() string {
	driver := strings.TrimSpace(strings.ToLower(config.Envs.MessageBus.Driver))
	if driver == "" {
		return "postgres"
	}

	return driver
}
