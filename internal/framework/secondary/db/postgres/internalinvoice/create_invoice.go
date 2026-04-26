package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) CreateInvoice(
	ctx context.Context,
	data coreentity.InternalInvoice,
) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:create_invoice:CreateInvoice",
	)
	defer span.End()

	number, err := r.nextNumber(ctx, "INV")
	if err != nil {
		return nil, err
	}
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_invoices (
			internal_order_id,
			invoice_number,
			status,
			currency_code,
			amount,
			amount_paid,
			amount_outstanding,
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
			?
		)
		RETURNING
			id,
			invoice_number,
			created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query),
		data.InternalOrderID,
		number,
		data.Status,
		data.CurrencyCode,
		data.Amount,
		data.AmountPaid,
		data.AmountOutstanding,
		data.DueAt,
		metadata,
	).Scan(&data.ID, &data.InvoiceNumber, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}
