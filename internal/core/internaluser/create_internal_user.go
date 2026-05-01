package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalUserCore) CreateInternalUser(
	ctx context.Context,
	data coreentity.InternalUser,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:CreateInternalUser")
	defer span.End()

	if err := c.authorizeManage(data.UserCtx); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidate(&data, true); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalUserByEmail(ctx, data.Email, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserEmailAlreadyExists)
	}

	hashedPassword, err := hashInternalPassword(data.Password)
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToHashPassword)
	}
	data.Password = hashedPassword

	return c.repo.CreateInternalUser(ctx, data)
}
