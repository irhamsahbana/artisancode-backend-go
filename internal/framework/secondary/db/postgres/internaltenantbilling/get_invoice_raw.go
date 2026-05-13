package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetInvoiceRaw(
	ctx context.Context,
	tenantID string,
	id string,
) (*coreentity.InternalTenantInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_invoice_raw:GetInvoiceRaw",
	)
	defer span.End()

	var item coreentity.InternalTenantInvoice
	query := `
		SELECT
			id,
			tenant_id::text,
			internal_billing_account_id,
			internal_tenant_subscription_id::text,
			internal_tenant_subscription_change_id::text,
			invoice_number,
			status,
			currency_code,
			amount::text,
			amount_paid::text,
			amount_outstanding::text,
			source_type,
			source_reference_id::text,
			target_subscription_state,
			due_at::text,
			paid_at::text,
			expired_at::text,
			created_at::text,
			updated_at::text
		FROM internal_tenant_invoices
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		id,
		tenantID,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", id).Str("tenant_id", tenantID).Msg("Failed to get tenant billing invoice raw")
		return nil, err
	}

	return &item, nil
}
