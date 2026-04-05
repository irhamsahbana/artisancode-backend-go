package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userCore) DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteUser")
	defer span.End()

	if filter.UserCtx.UserID == filter.ID {
		return errmsg.NewCustomErrors(400).SetMessage("You cannot delete your own account")
	}

	return c.repo.DeleteUser(ctx, filter)
}
