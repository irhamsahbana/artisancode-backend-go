package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) CreateInvoice(
	ctx context.Context,
	data coreentity.InternalTenantInvoice,
) (*coreentity.InternalTenantInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:create_invoice:CreateInvoice",
	)
	defer span.End()

	number, err := r.nextInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}

	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_tenant_invoices (
			tenant_id,
			internal_billing_account_id,
			internal_tenant_subscription_id,
			internal_tenant_subscription_change_id,
			invoice_number,
			status,
			currency_code,
			amount,
			amount_paid,
			amount_outstanding,
			source_type,
			source_reference_id,
			target_subscription_state,
			due_at,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id,
			invoice_number,
			created_at::text AS created_at
	`
	if err := r.exec(ctx).QueryRowxContext(
		ctx,
		r.exec(ctx).Rebind(query),
		data.TenantID,
		data.InternalBillingAccountID,
		data.InternalTenantSubscriptionID,
		data.InternalTenantSubscriptionChangeID,
		number,
		data.Status,
		data.CurrencyCode,
		data.Amount,
		data.AmountPaid,
		data.AmountOutstanding,
		data.SourceType,
		data.SourceReferenceID,
		data.TargetSubscriptionState,
		data.DueAt,
		metadata,
	).Scan(&data.ID, &data.InvoiceNumber, &data.CreatedAt); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create tenant billing invoice")
		return nil, err
	}

	return &data, nil
}
