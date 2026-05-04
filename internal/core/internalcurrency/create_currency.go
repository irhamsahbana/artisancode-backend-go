package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) CreateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) (*coreentity.InternalCurrency, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:create_currency:CreateInternalCurrency")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}
	if err := c.normalizeAndValidateCurrency(ctx, &data); err != nil {
		return nil, err
	}

	return c.repo.CreateInternalCurrency(ctx, data)
}
