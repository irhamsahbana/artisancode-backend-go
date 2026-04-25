package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalUserCore interface {
	Login(ctx context.Context, user coreentity.InternalUser) (*coreentity.AuthTokens, error)
	RefreshToken(ctx context.Context, user coreentity.InternalUser) (*coreentity.AuthTokens, error)
	GetInternalUsers(ctx context.Context, filter coreentity.InternalUserListFilter) ([]coreentity.InternalUser, int, error)
	GetInternalUser(ctx context.Context, filter coreentity.InternalUserFilter) (*coreentity.InternalUser, error)
	CreateInternalUser(ctx context.Context, data coreentity.InternalUser) (*coreentity.InternalUser, error)
	UpdateInternalUser(ctx context.Context, data coreentity.InternalUser) error
	DeleteInternalUser(ctx context.Context, filter coreentity.InternalUserDeleteFilter) error
}
