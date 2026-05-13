package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetInvoice(
	ctx context.Context,
	tenantID string,
	id string,
) (*coreentity.TenantBillingInvoiceDetail, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_invoice_by_id:GetInvoice",
	)
	defer span.End()

	var detail coreentity.TenantBillingInvoiceDetail
	query := `
		SELECT
			id,
			invoice_number,
			status,
			currency_code,
			amount::text,
			amount_paid::text,
			amount_outstanding::text,
			source_type,
			due_at::text,
			paid_at::text,
			expired_at::text,
			created_at::text
		FROM internal_tenant_invoices
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&detail,
		r.exec(ctx).Rebind(query),
		id,
		tenantID,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", id).Str("tenant_id", tenantID).Msg("Failed to get tenant billing invoice")
		return nil, err
	}

	attempts, err := r.GetPaymentAttemptsByInvoice(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	detail.PaymentAttempts = attempts

	return &detail, nil
}
