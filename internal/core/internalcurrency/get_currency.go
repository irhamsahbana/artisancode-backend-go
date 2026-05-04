package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) GetInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyFilter,
) (*coreentity.InternalCurrency, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:get_currency:GetInternalCurrency")
	defer span.End()

	return c.repo.GetInternalCurrency(ctx, filter)
}
