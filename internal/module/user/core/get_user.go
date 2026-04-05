package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *userCore) GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetUser")
	defer span.End()

	return c.repo.GetUser(ctx, filter)
}
