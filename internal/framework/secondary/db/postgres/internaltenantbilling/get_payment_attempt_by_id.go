package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetPaymentAttemptByID(
	ctx context.Context,
	tenantID string,
	id string,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_payment_attempt_by_id:GetPaymentAttemptByID",
	)
	defer span.End()

	var item coreentity.InternalTenantPaymentAttempt
	query := `
		SELECT
			id,
			tenant_id::text,
			internal_tenant_invoice_id,
			provider,
			payment_method_type,
			payment_channel_code,
			provider_reference,
			provider_request_id,
			provider_payment_url,
			status,
			requested_amount::text,
			paid_amount::text,
			expired_at::text,
			paid_at::text,
			failed_at::text,
			cancelled_at::text,
			raw_last_status,
			created_at::text,
			updated_at::text
		FROM internal_tenant_payment_attempts
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
		log.Ctx(ctx).Error().Err(err).Str("id", id).Str("tenant_id", tenantID).Msg("Failed to get tenant billing payment attempt by ID")
		return nil, err
	}

	return &item, nil
}
