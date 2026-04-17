package cmd

import (
	"codebase-app/internal/adapter"
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
)

const clearDataConfirmationToken = "DELETE_ALL_DATA"

var clearDataTables = []string{
	"attendance_logs",
	"export_jobs",
	"storage_file_links",
	"storage_files",
	"message_queue",
	"message_queue_dead_letters",
	"employees",
	"user_roles",
	"role_permissions",
	"work_locations",
	"work_shifts",
	"job_positions",
	"users",
	"roles",
	"permissions",
	"org_units",
	"tenants",
}

func RunClearData(cmd *flag.FlagSet, args []string) {
	confirm := cmd.String("confirm", "", "confirmation token required to clear all application data")

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing clear-data flags")
	}

	if *confirm != clearDataConfirmationToken {
		log.Fatal().
			Str("expected_confirm", clearDataConfirmationToken).
			Msg("Clear data aborted because confirmation token is invalid")
	}

	adapter.Adapters.Sync(
		adapter.WithPostgres(),
	)
	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
	}()

	if err := clearAllApplicationData(context.Background(), adapter.Adapters.Postgres); err != nil {
		log.Fatal().Err(err).Msg("Error while clearing application data")
	}

	log.Info().
		Strs("tables", clearDataTables).
		Msg("All application data cleared successfully")
}

func clearAllApplicationData(ctx context.Context, db rebindExecutor) error {
	query := `
		TRUNCATE TABLE
			attendance_logs,
			export_jobs,
			storage_file_links,
			storage_files,
			message_queue,
			employees,
			user_roles,
			role_permissions,
			work_locations,
			work_shifts,
			job_positions,
			users,
			roles,
			permissions,
			org_units,
			tenants
		RESTART IDENTITY CASCADE
	`

	_, err := db.ExecContext(ctx, db.Rebind(query))
	if err != nil {
		return fmt.Errorf("truncate application tables: %w", err)
	}

	return nil
}

type rebindExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Rebind(query string) string
}
