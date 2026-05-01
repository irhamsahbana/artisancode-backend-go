package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalUserCore) GetInternalUser(
	ctx context.Context,
	filter coreentity.InternalUserFilter,
) (*coreentity.InternalUser, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:GetInternalUser")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return nil, err
	}

	return c.repo.GetInternalUser(ctx, filter)
}
