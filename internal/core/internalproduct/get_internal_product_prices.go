package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalProductCore) GetInternalProductPrices(
	ctx context.Context,
	filter coreentity.InternalProductPriceListFilter,
) ([]coreentity.InternalProductPrice, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:GetInternalProductPrices",
	)
	defer span.End()

	if filter.CurrencyCode != "" {
		filter.CurrencyCode = strings.ToUpper(strings.TrimSpace(filter.CurrencyCode))
	}

	return c.repo.GetInternalProductPrices(ctx, filter)
}
