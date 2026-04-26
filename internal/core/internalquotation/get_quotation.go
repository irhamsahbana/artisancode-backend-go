package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalQuotationCore) GetQuotation(ctx context.Context, id string) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:get_quotation:GetQuotation")
	defer span.End()

	return c.quotationRepo.GetQuotation(ctx, userTenant(ctx), id)
}
