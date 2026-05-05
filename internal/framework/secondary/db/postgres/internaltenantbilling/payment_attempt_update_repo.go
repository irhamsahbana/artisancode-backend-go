package repository

import (
	"context"
	"encoding/json"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) UpdatePaymentAttemptStatus(
	ctx context.Context,
	id string,
	status string,
	metadata map[string]any,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:update_payment_attempt_status:UpdatePaymentAttemptStatus",
	)
	defer span.End()

	metaJSON, _ := json.Marshal(metadata)
	now := time.Now().UTC().Format(time.RFC3339)
	var paidAt *string
	var failedAt *string
	var cancelledAt *string

	switch status {
	case coreentity.PaymentAttemptStatusSucceeded:
		paidAt = &now
	case coreentity.PaymentAttemptStatusFailed:
		failedAt = &now
	case coreentity.PaymentAttemptStatusCancelled:
		cancelledAt = &now
	}

	query := `
		UPDATE internal_tenant_payment_attempts
		SET
			status = ?,
			paid_at = ?,
			failed_at = ?,
			cancelled_at = ?,
			metadata = nvl_merge_jsonb(metadata, ?),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		status,
		paidAt,
		failedAt,
		cancelledAt,
		metaJSON,
		id,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, id).Msg("Failed to update tenant billing payment attempt status")
		return err
	}

	return nil
}

func (r *internalTenantBillingRepo) UpdateInvoiceStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:update_invoice_status:UpdateInvoiceStatus",
	)
	defer span.End()

	now := time.Now().UTC().Format(time.RFC3339)
	var paidAt *string
	if status == coreentity.InvoiceStatusPaid {
		paidAt = &now
	}

	query := `
		UPDATE internal_tenant_invoices
		SET
			status = ?,
			paid_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		status,
		paidAt,
		id,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, id).Msg("Failed to update tenant billing invoice status")
		return err
	}

	return nil
}

func (r *internalTenantBillingRepo) GetInvoiceByNumber(
	ctx context.Context,
	number string,
) (*coreentity.InternalTenantInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_invoice_by_number:GetInvoiceByNumber",
	)
	defer span.End()

	var item coreentity.InternalTenantInvoice
	query := `
		SELECT
			id,
			tenant_id::text,
			internal_billing_account_id,
			internal_tenant_subscription_id::text,
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
			COALESCE(updated_at, created_at)::text AS updated_at
		FROM internal_tenant_invoices
		WHERE invoice_number = ?
			AND deleted_at IS NULL
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		number,
	); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalTenantBillingRepo) GetPaymentAttemptByProviderRef(
	ctx context.Context,
	providerReference string,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_payment_attempt_by_ref:GetPaymentAttemptByProviderRef",
	)
	defer span.End()

	var item coreentity.InternalTenantPaymentAttempt
	query := `
		SELECT
			id,
			tenant_id::text,
			internal_tenant_invoice_id,
			provider,
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
			COALESCE(updated_at, created_at)::text AS updated_at
		FROM internal_tenant_payment_attempts
		WHERE provider_reference = ?
			AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		providerReference,
	); err != nil {
		return nil, err
	}

	return &item, nil
}
