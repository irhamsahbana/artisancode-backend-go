package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserRepository interface {
	ExistsTenantByCode(ctx context.Context, code string) (bool, error)
	InsertTenant(ctx context.Context, tenant coreentity.Tenant) (string, error)
	FindActiveUserByEmailAndTenant(ctx context.Context, email, tenantID string) (*coreentity.User, error)
	ExistsActiveUserByEmail(ctx context.Context, email string) (bool, error)
	GetRoleByName(ctx context.Context, roleName, tenantID string) (*coreentity.Role, error)
	InsertUser(ctx context.Context, user coreentity.User) (string, error)

	InitializeTenant(ctx context.Context, tenantID string, companyName string) (string, error)
}