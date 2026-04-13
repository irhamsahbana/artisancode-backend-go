package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *userCore) UpdateUser(ctx context.Context, data coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:update_user:UpdateUser")
	defer span.End()

	emailExists, err := c.repo.ExistsActiveUserByEmailExcludeUser(ctx, data.Email, data.ID)
	if err != nil {
		return err
	}
	if emailExists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"email":   data.Email,
			"user_id": data.ID,
		}).Msg("Email already registered")
		return errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
	}

	if data.Password != "" {
		hashedPassword, err := hashPassword(data.Password)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to hash password")
			return errmsg.NewCustomErrors(500).SetMessage("Failed to update user")
		}
		data.Password = hashedPassword
	}

	return c.repo.UpdateUser(ctx, data)
}
