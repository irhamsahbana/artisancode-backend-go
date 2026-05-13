package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *userCore) CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:create_user:CreateUser")
	defer span.End()

	if c.billingCore != nil {
		count, err := c.repo.CountUsersByTenant(ctx, data.TenantID)
		if err != nil {
			return nil, err
		}
		if err := c.billingCore.CheckUsageLimit(ctx, coreentity.ResourceTypeUsers, count); err != nil {
			return nil, err
		}
	}

	emailExists, err := c.repo.ExistsActiveUserByEmailAndTenant(ctx, data.Email, data.TenantID)
	if err != nil {
		return nil, err
	}
	if emailExists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"email": data.Email,
		}).Msg("Email already registered")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageEmailIsAlreadyRegistered)
	}

	hashedPassword, err := hashPassword(data.Password)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg(errmsg.MessageFailedToHashPassword)
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageFailedToCreateUser)
	}

	data.Password = hashedPassword

	return c.repo.CreateUser(ctx, data)
}
