package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetPaymentAttemptsByInvoice(
	ctx context.Context,
	tenantID string,
	invoiceID string,
) ([]coreentity.TenantBillingPaymentAttemptView, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_payment_attempts:GetPaymentAttemptsByInvoice",
	)
	defer span.End()

	items := make([]coreentity.TenantBillingPaymentAttemptView, 0)
	query := `
		SELECT
			id,
			provider,
			payment_method_type,
			payment_channel_code,
			provider_reference,
			provider_payment_url,
			status,
			requested_amount::text,
			paid_amount::text,
			expired_at::text,
			paid_at::text,
			failed_at::text,
			created_at::text
		FROM internal_tenant_payment_attempts
		WHERE internal_tenant_invoice_id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	if err := r.exec(ctx).SelectContext(
		ctx,
		&items,
		r.exec(ctx).Rebind(query),
		invoiceID,
		tenantID,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("invoice_id", invoiceID).Str("tenant_id", tenantID).Msg("Failed to get tenant billing payment attempts")
		return nil, err
	}

	return items, nil
}
