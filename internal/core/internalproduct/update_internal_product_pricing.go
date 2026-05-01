package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) UpdateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalproduct:core:UpdateInternalProductPricing",
	)
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if _, err := c.repo.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: data.InternalProductID}); err != nil {
		return err
	}

	if err := c.normalizeAndValidatePricing(ctx, &data); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalProductPricingByCode(ctx, data.InternalProductID, data.Code, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductPricingCodeAlreadyExists)
	}

	return c.repo.UpdateInternalProductPricing(ctx, data)
}
