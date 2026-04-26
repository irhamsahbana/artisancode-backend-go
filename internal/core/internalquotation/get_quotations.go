package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalQuotationCore) GetQuotations(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalQuotation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:get_quotations:GetQuotations")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = userTenant(ctx)
	}
	return c.quotationRepo.GetQuotations(ctx, filter)
}
