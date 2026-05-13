package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) CancelInvoice(
	ctx context.Context,
	tenantID string,
	id string,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:cancel_invoice:CancelInvoice",
	)
	defer span.End()

	query := `
		UPDATE internal_tenant_invoices
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
			AND status = ?
	`
	result, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		coreentity.InternalTenantInvoiceStatusCancelled,
		id,
		tenantID,
		coreentity.InternalTenantInvoiceStatusPending,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", id).Msg("Failed to cancel tenant billing invoice")
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil
	}

	return nil
}
