package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) GetProviderCurrencies(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyListFilter,
) ([]coreentity.InternalPaymentProviderCurrency, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:get_provider_currencies:GetProviderCurrencies")
	defer span.End()

	filter.Provider = strings.ToLower(strings.TrimSpace(filter.Provider))
	return c.repo.GetProviderCurrencies(ctx, filter)
}
