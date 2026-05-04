package repository

import (
	"context"
	"fmt"
	"time"

	"codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type internalTenantBillingRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalTenantBillingRepository = &internalTenantBillingRepo{}

func NewInternalTenantBillingRepository(cfg Config) portsRepo.InternalTenantBillingRepository {
	return &internalTenantBillingRepo{db: cfg.DB}
}

func (r *internalTenantBillingRepo) exec(ctx context.Context) transaction.SQLExecutor {
	return transaction.ExecutorFromContext(ctx, r.db)
}

func (r *internalTenantBillingRepo) nextInvoiceNumber(ctx context.Context) (string, error) {
	exec := r.exec(ctx)
	today := time.Now().Format("20060102")
	like := fmt.Sprintf("TINV-%s-%%", today)
	query := `SELECT COUNT(*) + 1 FROM internal_tenant_invoices WHERE invoice_number LIKE ?`

	var seq int
	if err := exec.GetContext(ctx, &seq, exec.Rebind(query), like); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get next tenant billing invoice number")
		return "", err
	}

	return fmt.Sprintf("TINV-%s-%04d", today, seq), nil
}
