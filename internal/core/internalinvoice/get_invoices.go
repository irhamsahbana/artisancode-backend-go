package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalInvoiceCore) GetInvoices(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalInvoice, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:get_invoices:GetInvoices")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = common.GetUserContext(ctx).TenantID
	}
	return c.invoiceRepo.GetInvoices(ctx, filter)
}
