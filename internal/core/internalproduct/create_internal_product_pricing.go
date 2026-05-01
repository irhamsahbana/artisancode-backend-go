package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) CreateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:CreateInternalProductPricing",
	)
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if _, err := c.repo.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: data.InternalProductID}); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidatePricing(ctx, &data); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalProductPricingByCode(ctx, data.InternalProductID, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductPricingCodeAlreadyExists)
	}

	return c.repo.CreateInternalProductPricing(ctx, data)
}
