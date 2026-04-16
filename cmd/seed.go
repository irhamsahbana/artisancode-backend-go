package cmd

import (
	"codebase-app/db/seeds"
	"codebase-app/internal/adapter"
	"flag"

	"github.com/rs/zerolog/log"
)

// go run cmd/bin/main.go seed -table=ingredients_embeddings
func RunSeed(cmd *flag.FlagSet, args []string) {
	table, total := seeds.BindSeedFlags(cmd)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	adapter.Adapters.Sync(
		adapter.WithPostgres(),
	)
	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
	}()

	if err := seeds.Execute(adapter.Adapters.Postgres, *table, *total); err != nil {
		log.Fatal().Err(err).Msg("Error while running seed")
	}
}
