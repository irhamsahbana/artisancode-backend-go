package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalUserCore) DeleteInternalUser(ctx context.Context, filter coreentity.InternalUserDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:core:DeleteInternalUser")
	defer span.End()

	if err := c.authorizeManage(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalUser(ctx, filter)
}
