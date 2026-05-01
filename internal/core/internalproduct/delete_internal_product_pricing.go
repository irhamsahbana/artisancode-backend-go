package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) DeleteInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:DeleteInternalProductPricing",
	)
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProductPricing(ctx, filter)
}
