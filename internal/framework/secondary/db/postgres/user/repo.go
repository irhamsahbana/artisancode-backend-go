package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	repository "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

var _ repository.UserRepository = &userRepo{}

type userRepo struct {
	db *sqlx.DB
}

func (r *userRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

type Config struct {
	DB *sqlx.DB
}

func NewUserRepository(cfg Config) repository.UserRepository {
	return &userRepo{
		db: cfg.DB,
	}
}
