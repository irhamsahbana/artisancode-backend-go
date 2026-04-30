package transaction

import (
	"context"
	"database/sql"

	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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

	if fallback == nil {
		log.Ctx(ctx).Debug().Msg("No SQLExecutor found in context and fallback is nil")
		return nil
	}

	log.Ctx(ctx).Debug().Msg("No SQLExecutor found in context, using fallback")
	return fallback
}

func (t *transactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) (err error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:transaction:transactor:WithinTransaction",
	)
	defer span.End()
	log.Ctx(ctx).Debug().Msg("Starting transaction")

	if existing := ExecutorFromContext(ctx, nil); existing != nil {
		log.Ctx(ctx).Debug().Msg("Already within a transaction, executing function directly")
		return fn(ctx)
	}

	log.Ctx(ctx).Debug().Msg("Beginning new transaction")
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to begin transaction")
		return err
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Ctx(ctx).Debug().Msg("Recovering from panic")
			t.rollbackTransaction(ctx, tx)
			panic(recovered)
		}

		if err != nil {
			log.Ctx(ctx).Debug().Msg("Error occurred, rolling back transaction")
			t.rollbackTransaction(ctx, tx)
			return
		}

		err = t.commitTransaction(ctx, tx)
	}()

	log.Ctx(ctx).Debug().Msg("Executing function within transaction")
	return fn(WithExecutor(ctx, tx))
}

func (t *transactor) rollbackTransaction(ctx context.Context, tx *sqlx.Tx) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:transaction:transactor:rollbackTransaction",
	)
	defer span.End()

	log.Ctx(ctx).Debug().Msg("Rolling back transaction")
	err := tx.Rollback()
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("Failed to rollback transaction")
	}
}

func (t *transactor) commitTransaction(ctx context.Context, tx *sqlx.Tx) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:transaction:transactor:commitTransaction",
	)
	defer span.End()

	log.Ctx(ctx).Debug().Msg("Committing transaction")
	err := tx.Commit()
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("Failed to commit transaction")
	}

	return err
}
