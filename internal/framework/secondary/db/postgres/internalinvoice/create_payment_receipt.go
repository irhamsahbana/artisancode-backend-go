package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) CreatePaymentReceipt(
	ctx context.Context,
	data coreentity.InternalPaymentReceipt,
) (*coreentity.InternalPaymentReceipt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:create_payment_receipt:CreatePaymentReceipt",
	)
	defer span.End()

	number, err := r.nextNumber(ctx, "RCP")
	if err != nil {
		return nil, err
	}
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_payment_receipts (
			internal_invoice_id,
			internal_payment_attempt_id,
			receipt_number,
			status,
			amount_received,
			currency_code,
			received_at,
			verified_at,
			verified_by_user_id,
			source_type,
			reference_number,
			notes,
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
			?
		)
		RETURNING
			id,
			receipt_number,
			created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query),
		data.InternalInvoiceID,
		data.InternalPaymentAttemptID,
		number,
		data.Status,
		data.AmountReceived,
		data.CurrencyCode,
		data.ReceivedAt,
		data.VerifiedAt,
		data.VerifiedByUserID,
		data.SourceType,
		data.ReferenceNumber,
		data.Notes,
		metadata,
	).Scan(&data.ID, &data.ReceiptNumber, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}
