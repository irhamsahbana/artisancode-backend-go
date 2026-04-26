package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) GetInvoiceByOrderID(
	ctx context.Context,
	tenantID, orderID string,
) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_invoice_by_order_id:GetInvoiceByOrderID",
	)
	defer span.End()

	return r.getInvoiceByWhere(ctx, tenantID, "i.internal_order_id = ?", orderID)
}
