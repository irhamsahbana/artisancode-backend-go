package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) GetInternalCurrencies(
	ctx context.Context,
	filter coreentity.InternalCurrencyListFilter,
) ([]coreentity.InternalCurrency, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:get_currencies:GetInternalCurrencies")
	defer span.End()

	return c.repo.GetInternalCurrencies(ctx, filter)
}
