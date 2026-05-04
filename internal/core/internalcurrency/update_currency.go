package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) UpdateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:update_currency:UpdateInternalCurrency")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}
	if err := c.normalizeAndValidateCurrency(ctx, &data); err != nil {
		return err
	}

	return c.repo.UpdateInternalCurrency(ctx, data)
}
