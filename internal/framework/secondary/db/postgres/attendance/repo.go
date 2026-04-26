package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

type attendanceRepo struct {
	db *sqlx.DB
}

func (r *attendanceRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

var _ portsRepo.AttendanceRepository = &attendanceRepo{}

type Config struct {
	DB *sqlx.DB
}

func NewAttendanceRepository(cfg Config) portsRepo.AttendanceRepository {
	return &attendanceRepo{db: cfg.DB}
}
