package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetInvoices(
	ctx context.Context,
	filter coreentity.TenantBillingInvoiceListFilter,
) ([]coreentity.TenantBillingInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_invoices:GetInvoices")
	defer span.End()

	if c.billingRepo == nil {
		return []coreentity.TenantBillingInvoice{}, nil
	}

	return c.billingRepo.GetInvoices(ctx, filter)
}
