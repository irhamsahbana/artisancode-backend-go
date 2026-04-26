package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) GetLatestPaymentAttempt(
	ctx context.Context,
	tenantID string,
	invoiceID string,
) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_latest_payment_attempt:GetLatestPaymentAttempt",
	)
	defer span.End()

	items, err := r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.internal_invoice_id = ?", invoiceID, 1)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}
