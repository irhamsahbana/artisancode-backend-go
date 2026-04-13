package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *userCore) GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:get_users:GetUsers")
	defer span.End()

	return c.repo.GetUsers(ctx, filter)
}
