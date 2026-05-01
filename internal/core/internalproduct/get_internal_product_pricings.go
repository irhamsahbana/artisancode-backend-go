package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) GetInternalProductPricings(
	ctx context.Context,
	filter coreentity.InternalProductPricingListFilter,
) ([]coreentity.InternalProductPricing, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:GetInternalProductPricings",
	)
	defer span.End()

	return c.repo.GetInternalProductPricings(ctx, filter)
}

func (c *internalProductCore) GetInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingFilter,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:GetInternalProductPricing",
	)
	defer span.End()

	return c.repo.GetInternalProductPricing(ctx, filter)
}
