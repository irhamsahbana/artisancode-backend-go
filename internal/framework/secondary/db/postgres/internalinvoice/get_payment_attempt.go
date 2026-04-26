package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *internalInvoiceRepo) GetPaymentAttempt(
	ctx context.Context,
	tenantID, id string,
) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_payment_attempt:GetPaymentAttempt",
	)
	defer span.End()

	items, err := r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.id = ?", id, 1)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Payment attempt not found")
	}
	return &items[0], nil
}
