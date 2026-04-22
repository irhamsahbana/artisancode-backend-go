package adapter

import (
	"fmt"
	"time"

	"codebase-app/internal/framework/secondary/publisher/postgres"
	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
)

func WithMessageSubscriptionManager() Option {
	return func(a *Adapter) {
		if a.Postgres == nil {
			err := fmt.Errorf("Postgres adapter must be initialized before Message Subscription Manager")
			log.Fatal().Err(err).Msg("Failed to initialize message subscription manager")
		}

		a.MessageSubscriptionManager = postgres.NewSubscriptionManager(a.Postgres, postgres.Config{
			PollInterval:      time.Duration(config.Envs.MessageBus.PostgresPollIntervalMS) * time.Millisecond,
			BatchSize:         config.Envs.MessageBus.PostgresBatchSize,
			RetryDelay:        time.Duration(config.Envs.MessageBus.RetryDelaySeconds) * time.Second,
			MaxAttempts:       config.Envs.MessageBus.MaxAttempts,
			ProcessingTimeout: time.Duration(config.Envs.MessageBus.ProcessingTimeoutSecs) * time.Second,
		})

		log.Info().Msg("Message subscription manager connected with driver: postgres")
	}
}
