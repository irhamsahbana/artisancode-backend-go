package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *userCore) Logout(ctx context.Context, user coreentity.User) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:user:logout:Logout")
	defer span.End()

	if user.RefreshToken == "" {
		return nil
	}

	c.tokenCache.DeleteRefreshToken(ctx, user.RefreshToken)

	return nil
}
