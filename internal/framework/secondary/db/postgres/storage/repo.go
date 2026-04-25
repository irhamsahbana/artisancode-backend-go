package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

type storageRepo struct {
	db *sqlx.DB
}

func (r *storageRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

var _ portsRepo.StorageRepository = &storageRepo{}

type StorageRepositoryConfig struct {
	DB *sqlx.DB
}

func NewStorageRepository(cfg StorageRepositoryConfig) portsRepo.StorageRepository {
	return &storageRepo{db: cfg.DB}
}
