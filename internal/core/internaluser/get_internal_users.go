package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalUserCore) GetInternalUsers(
	ctx context.Context,
	filter coreentity.InternalUserListFilter,
) ([]coreentity.InternalUser, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:GetInternalUsers")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return nil, 0, err
	}

	return c.repo.GetInternalUsers(ctx, filter)
}
