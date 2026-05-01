package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalProductCore) CreateInternalProduct(
	ctx context.Context,
	data coreentity.InternalProduct,
) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:CreateInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidateProduct(ctx, &data); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalProductByCode(ctx, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalProductCodeAlreadyExists)
	}

	return c.repo.CreateInternalProduct(ctx, data)
}
