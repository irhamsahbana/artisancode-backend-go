package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *userRepo) FindActiveUserByEmailAndTenant(
	ctx context.Context,
	email, tenantCode string,
) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:login:FindActiveUserByEmailAndTenant",
	)
	defer span.End()

	// Get user base data
	user, err := r.GetUserByEmailAndTenant(ctx, email, tenantCode)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	// Get user roles
	roleNames, err := r.GetUserRolesByUserID(ctx, user.ID)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	// Attach roles to user object
	user.RoleNames = roleNames

	return user, nil
}
