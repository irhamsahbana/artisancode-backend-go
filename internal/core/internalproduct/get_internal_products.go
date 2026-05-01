package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) GetInternalProducts(
	ctx context.Context,
	filter coreentity.InternalProductListFilter,
) ([]coreentity.InternalProduct, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProducts")
	defer span.End()

	return c.repo.GetInternalProducts(ctx, filter)
}

func (c *internalProductCore) GetInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductFilter,
) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProduct")
	defer span.End()

	return c.repo.GetInternalProduct(ctx, filter)
}
