package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *userRepo) CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:create_user:CreateUser")
	defer span.End()

	userID, err := r.InsertUser(ctx, data)
	if err != nil {
		return nil, err
	}

	return r.GetUser(ctx, coreentity.User{
		ID:       userID,
		TenantID: data.TenantID,
	})
}
