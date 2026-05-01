package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalUserCore) Logout(ctx context.Context, user coreentity.InternalUser) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaluser:logout:Logout")
	defer span.End()

	if user.RefreshToken == "" {
		return nil
	}

	c.tokenCache.DeleteRefreshToken(ctx, user.RefreshToken)

	return nil
}
