package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserCore interface {
	Login(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	RefreshToken(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	RegisterOwner(ctx context.Context, user coreentity.User, tenant coreentity.Tenant) (*coreentity.AuthTokens, error)
}
