package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetAllInvoices(
	ctx context.Context,
	filter coreentity.InternalBillingInvoiceListFilter,
) ([]coreentity.InternalTenantInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_all_invoices:GetAllInvoices")
	defer span.End()

	return c.billingRepo.GetAllInvoices(ctx, filter)
}
