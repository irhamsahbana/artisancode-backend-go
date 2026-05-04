package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) DeleteProviderCurrency(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyFilter,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:delete_provider_currency:DeleteProviderCurrency")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	filter.Provider = strings.ToLower(strings.TrimSpace(filter.Provider))
	filter.CurrencyCode = strings.ToUpper(strings.TrimSpace(filter.CurrencyCode))
	return c.repo.DeleteProviderCurrency(ctx, filter)
}
