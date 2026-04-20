package adapter

import (
	"fmt"
	"strings"

	"codebase-app/internal/framework/secondary/publisher/nats"
	"codebase-app/internal/framework/secondary/publisher/postgres"
	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
)

func WithMessagePublisher() Option {
	return func(a *Adapter) {
		var err error

		switch normalizedMessageBusDriver() {
		case "postgres":
			if a.Postgres == nil {
				err = fmt.Errorf("Postgres adapter must be initialized before Message Publisher")
				log.Fatal().Err(err).Msg("Failed to initialize message publisher")
			}
			a.MessagePublisher = postgres.NewPublisher(a.Postgres)
		case "nats":
			a.MessagePublisher, err = nats.NewPublisher(config.Envs.EmailVerificationQueueNats.NatsURL)
		default:
			err = fmt.Errorf("unsupported message bus driver: %s", config.Envs.MessageBus.Driver)
		}

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize message publisher")
		}

		log.Info().Msgf("Message publisher initialized with driver: %s", normalizedMessageBusDriver())
	}
}

func normalizedMessageBusDriver() string {
	driver := strings.TrimSpace(strings.ToLower(config.Envs.MessageBus.Driver))
	if driver == "" {
		return "postgres"
	}

	return driver
}
