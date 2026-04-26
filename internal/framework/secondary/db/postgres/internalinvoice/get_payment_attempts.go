package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) GetPaymentAttempts(
	ctx context.Context,
	tenantID string,
	invoiceID string,
) ([]coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_payment_attempts:GetPaymentAttempts",
	)
	defer span.End()

	return r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.internal_invoice_id = ?", invoiceID, 100)
}
