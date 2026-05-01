package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:UpdateInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if err := c.normalizeAndValidateProduct(ctx, &data); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalProductByCode(ctx, data.Code, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductCodeAlreadyExists)
	}

	return c.repo.UpdateInternalProduct(ctx, data)
}
