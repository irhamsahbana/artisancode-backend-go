package adapter

import (
	"fmt"

	"codebase-app/internal/framework/secondary/publisher/postgres"

	"github.com/rs/zerolog/log"
)

func WithMessagePublisher() Option {
	return func(a *Adapter) {
		if a.Postgres == nil {
			err := fmt.Errorf("Postgres adapter must be initialized before Message Publisher")
			log.Fatal().Err(err).Msg("Failed to initialize message publisher")
		}

		a.MessagePublisher = postgres.NewPublisher(a.Postgres)

		log.Info().Msg("Message publisher initialized with driver: postgres")
	}
}
