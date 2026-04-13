package adapter

import (
	secondarypostgres "codebase-app/internal/framework/secondary/db/postgres"
	"codebase-app/internal/infrastructure/config"

	"github.com/rs/zerolog/log"
)

func WithPostgres() Option {
	return func(a *Adapter) {
		db, err := secondarypostgres.New(secondarypostgres.Config{
			Username:        config.Envs.Postgres.Username,
			Password:        config.Envs.Postgres.Password,
			Database:        config.Envs.Postgres.Database,
			Host:            config.Envs.Postgres.Host,
			Port:            config.Envs.Postgres.Port,
			SSLMode:         config.Envs.Postgres.SslMode,
			MaxOpenConns:    config.Envs.DB.MaxOpenCons,
			MaxIdleConns:    config.Envs.DB.MaxIdleCons,
			ConnMaxLifetime: config.Envs.DB.ConnMaxLifetime,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Postgres")
		}

		if err := db.Ping(); err != nil {
			log.Fatal().Err(err).Msg("Error pinging Postgres")
		}

		a.Postgres = db
		log.Info().Msg("Postgres connected")
	}
}
