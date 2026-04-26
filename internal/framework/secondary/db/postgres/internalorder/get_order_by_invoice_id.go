package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalOrderRepo) GetOrderByInvoiceID(
	ctx context.Context,
	tenantID, invoiceID string,
) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalorder:get_order_by_invoice_id:GetOrderByInvoiceID",
	)
	defer span.End()

	return r.getOrderByWhere(ctx, tenantID, "i.id = ?", invoiceID)
}
