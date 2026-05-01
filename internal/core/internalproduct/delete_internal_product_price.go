package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) DeleteInternalProductPrice(
	ctx context.Context,
	filter coreentity.InternalProductPriceDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:DeleteInternalProductPrice",
	)
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProductPrice(ctx, filter)
}
