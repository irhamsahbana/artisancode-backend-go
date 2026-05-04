package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalCurrencyCore) UpsertProviderCurrency(
	ctx context.Context,
	data coreentity.InternalPaymentProviderCurrency,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:upsert_provider_currency:UpsertProviderCurrency")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}
	if err := c.normalizeAndValidateProviderCurrency(ctx, &data); err != nil {
		return err
	}

	currency, err := c.repo.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{
		Code: data.CurrencyCode,
	})
	if err != nil {
		return err
	}
	if data.IsActive && !currency.IsActive {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyIsNotActive)
	}

	return c.repo.UpsertProviderCurrency(ctx, data)
}
