package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) CreatePaymentAttempt(
	ctx context.Context,
	data coreentity.InternalPaymentAttempt,
) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:create_payment_attempt:CreatePaymentAttempt",
	)
	defer span.End()

	payload, _ := json.Marshal(data.ProviderPayloadSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_payment_attempts (
			internal_invoice_id,
			provider,
			payment_method_type,
			payment_channel_code,
			provider_reference,
			provider_request_id,
			provider_payment_url,
			provider_payload_snapshot,
			status,
			requested_amount,
			paid_amount,
			expired_at,
			raw_last_status,
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
			?
		)
		RETURNING
			id,
			created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query),
		data.InternalInvoiceID,
		data.Provider,
		data.PaymentMethodType,
		data.PaymentChannelCode,
		data.ProviderReference,
		data.ProviderRequestID,
		data.ProviderPaymentURL,
		payload,
		data.Status,
		data.RequestedAmount,
		data.PaidAmount,
		data.ExpiredAt,
		data.RawLastStatus,
		metadata,
	).Scan(&data.ID, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}
