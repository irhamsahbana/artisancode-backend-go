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

type EmployeeRepositoryConfig struct {
	DB *sqlx.DB
}

func NewEmployeeRepository(cfg EmployeeRepositoryConfig) portsRepo.EmployeeRepository {
	return &employeeRepo{db: cfg.DB}
}
