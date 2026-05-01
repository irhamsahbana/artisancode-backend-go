package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) UpdateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:UpdateInternalProductPrice",
	)
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if _, err := c.repo.GetInternalProductPricing(ctx, coreentity.InternalProductPricingFilter{
		ID: data.InternalProductPricingID,
	}); err != nil {
		return err
	}

	if err := c.normalizeAndValidatePrice(ctx, &data); err != nil {
		return err
	}

	overlap, err := c.repo.ExistsOverlappingInternalProductPrice(ctx, coreentity.InternalProductPriceOverlapFilter{
		InternalProductPricingID: data.InternalProductPricingID,
		CurrencyCode:             data.CurrencyCode,
		StartedAt:                data.StartedAt,
		EndedAt:                  data.EndedAt,
		ExcludeID:                data.ID,
	})
	if err != nil {
		return err
	}
	if overlap {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePricePeriodOverlapsWithAnExistingActivePrice)
	}

	return c.repo.UpdateInternalProductPrice(ctx, data)
}
