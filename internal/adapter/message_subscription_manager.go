package adapter

import (
	"fmt"
	"time"

	"codebase-app/internal/framework/secondary/publisher/nats"
	"codebase-app/internal/framework/secondary/publisher/postgres"
	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
)

func WithMessageSubscriptionManager() Option {
	return func(a *Adapter) {
		var err error

		switch normalizedMessageBusDriver() {
		case "postgres":
			if a.Postgres == nil {
				err = fmt.Errorf("Postgres adapter must be initialized before Message Subscription Manager")
				log.Fatal().Err(err).Msg("Failed to initialize message subscription manager")
			}
			a.MessageSubscriptionManager = postgres.NewSubscriptionManager(a.Postgres, postgres.Config{
				PollInterval: time.Duration(config.Envs.MessageBus.PostgresPollIntervalMS) * time.Millisecond,
				BatchSize:    config.Envs.MessageBus.PostgresBatchSize,
				RetryDelay:   time.Duration(config.Envs.MessageBus.RetryDelaySeconds) * time.Second,
				MaxAttempts:  config.Envs.MessageBus.MaxAttempts,
			})
		case "nats":
			a.MessageSubscriptionManager, err = nats.NewSubscriptionManager(config.Envs.EmailVerificationQueueNats.NatsURL)
		default:
			err = fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
		}

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize message subscription manager")
		}

		log.Info().Msgf("Message subscription manager connected with driver: %s", normalizedMessageBusDriver())
	}
}
