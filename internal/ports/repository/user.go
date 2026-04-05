package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserRepository interface {
	ExistsTenantByCode(ctx context.Context, code string) (bool, error)
	InsertTenant(ctx context.Context, tenant coreentity.Tenant) (string, error)
	FindActiveUserByEmailAndTenant(ctx context.Context, email, tenantCode string) (*coreentity.User, error)
	ExistsActiveUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsActiveUserByEmailExcludeUser(ctx context.Context, email, excludeUserID string) (bool, error)
	GetRoleByName(ctx context.Context, roleName, tenantID string) (*coreentity.Role, error)
	InsertUser(ctx context.Context, user coreentity.User) (string, error)
	UpdateUserEmail(ctx context.Context, userID, tenantID, email string) error
	UpdateUserPassword(ctx context.Context, userID, tenantID, hashedPassword string) error
	GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error)
	GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error)
	CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error)
	UpdateUser(ctx context.Context, data coreentity.User) error
	DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error

	InitializeTenant(ctx context.Context, tenantID string, companyName string) (string, error)
}
