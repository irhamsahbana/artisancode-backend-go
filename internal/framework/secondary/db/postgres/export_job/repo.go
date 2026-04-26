package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

type exportJobRepo struct {
	db *sqlx.DB
}

func (r *exportJobRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

var _ portsRepo.ExportJobRepository = &exportJobRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewExportJobRepository(cfg Config) portsRepo.ExportJobRepository {
	return &exportJobRepo{db: cfg.DB}
}
