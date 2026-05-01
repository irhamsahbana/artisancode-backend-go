package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) CreateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) (*coreentity.InternalProductPrice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:CreateInternalProductPrice",
	)
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if _, err := c.repo.GetInternalProductPricing(ctx, coreentity.InternalProductPricingFilter{
		ID: data.InternalProductPricingID,
	}); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidatePrice(ctx, &data); err != nil {
		return nil, err
	}

	overlap, err := c.repo.ExistsOverlappingInternalProductPrice(ctx, coreentity.InternalProductPriceOverlapFilter{
		InternalProductPricingID: data.InternalProductPricingID,
		CurrencyCode:             data.CurrencyCode,
		StartedAt:                data.StartedAt,
		EndedAt:                  data.EndedAt,
	})
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePricePeriodOverlapsWithAnExistingActivePrice)
	}

	return c.repo.CreateInternalProductPrice(ctx, data)
}
