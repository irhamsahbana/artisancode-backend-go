package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

type employeeRepo struct {
	db *sqlx.DB
}

func (r *employeeRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

var _ portsRepo.EmployeeRepository = &employeeRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewEmployeeRepository(cfg Config) portsRepo.EmployeeRepository {
	return &employeeRepo{db: cfg.DB}
}
