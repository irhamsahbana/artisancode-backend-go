package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserCore interface {
	Login(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	RefreshToken(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	RegisterOwner(ctx context.Context, user coreentity.User, tenant coreentity.Tenant) (*coreentity.AuthTokens, error)
	GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error)
	GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error)
	CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error)
	UpdateUser(ctx context.Context, data coreentity.User) error
	DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error
}
