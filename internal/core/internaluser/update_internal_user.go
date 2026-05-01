package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalUserCore) UpdateInternalUser(ctx context.Context, data coreentity.InternalUser) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:UpdateInternalUser")
	defer span.End()

	if err := c.authorizeManage(data.UserCtx); err != nil {
		return err
	}

	if err := c.normalizeAndValidate(&data, false); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalUserByEmail(ctx, data.Email, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserEmailAlreadyExists)
	}

	if data.Password != "" {
		hashedPassword, err := hashInternalPassword(data.Password)
		if err != nil {
			return errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToHashPassword)
		}
		data.Password = hashedPassword
	}

	return c.repo.UpdateInternalUser(ctx, data)
}
