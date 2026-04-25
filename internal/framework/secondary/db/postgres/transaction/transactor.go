package transaction

import (
	"context"
	"database/sql"

	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const executorContextKey contextKey = "postgres-executor"

type SQLExecutor interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Rebind(string) string
}

type transactor struct {
	db *sqlx.DB
}

var _ portsRepo.Transactor = &transactor{}

func NewTransactor(db *sqlx.DB) portsRepo.Transactor {
	return &transactor{db: db}
}

func WithExecutor(ctx context.Context, exec SQLExecutor) context.Context {
	if exec == nil {
		return ctx
	}

	return context.WithValue(ctx, executorContextKey, exec)
}

func ExecutorFromContext(ctx context.Context, fallback *sqlx.DB) SQLExecutor {
	if ctx != nil {
		if exec, ok := ctx.Value(executorContextKey).(SQLExecutor); ok && exec != nil {
			return exec
		}
	}

	return fallback
}

func (t *transactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) (err error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:transaction:transactor:WithinTransaction")
	defer span.End()

	if existing := ExecutorFromContext(ctx, nil); existing != nil {
		return fn(ctx)
	}

	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			_ = tx.Rollback()
			panic(recovered)
		}

		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return fn(WithExecutor(ctx, tx))
}
