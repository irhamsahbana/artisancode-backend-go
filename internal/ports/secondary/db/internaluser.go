package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalUserRepository interface {
	FindActiveInternalUserByEmail(ctx context.Context, email string) (*coreentity.InternalUser, error)
	GetInternalUser(ctx context.Context, filter coreentity.InternalUserFilter) (*coreentity.InternalUser, error)
	GetInternalUsers(ctx context.Context, filter coreentity.InternalUserListFilter) ([]coreentity.InternalUser, int, error)
	CreateInternalUser(ctx context.Context, data coreentity.InternalUser) (*coreentity.InternalUser, error)
	UpdateInternalUser(ctx context.Context, data coreentity.InternalUser) error
	DeleteInternalUser(ctx context.Context, filter coreentity.InternalUserDeleteFilter) error
	UpdateInternalUserLastLogin(ctx context.Context, id string) error
	ExistsInternalUserByEmail(ctx context.Context, email, excludeID string) (bool, error)
}
