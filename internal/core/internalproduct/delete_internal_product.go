package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) DeleteInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:DeleteInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProduct(ctx, filter)
}
