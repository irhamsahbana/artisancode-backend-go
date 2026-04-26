package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) GetInvoice(
	ctx context.Context,
	tenantID, id string,
) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_invoice:GetInvoice",
	)
	defer span.End()

	return r.getInvoiceByWhere(ctx, tenantID, "i.id = ?", id)
}
