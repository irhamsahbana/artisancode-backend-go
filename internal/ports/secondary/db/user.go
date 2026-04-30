package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserRepository interface {
	ExistsTenantByCode(ctx context.Context, code string) (bool, error)
	InsertTenant(ctx context.Context, tenant coreentity.Tenant) (string, error)
	FindActiveUserByEmailAndTenant(ctx context.Context, email, tenantCode string) (*coreentity.User, error)
	FindActiveUsersByEmail(ctx context.Context, email string) ([]coreentity.User, error)
	FindActiveUserByEmailAndTenantID(ctx context.Context, email, tenantID string) (*coreentity.User, error)
	FindActiveUserByIDAndTenant(ctx context.Context, userID, tenantID string) (*coreentity.User, error)
	FindAuthIdentityByProviderSubject(ctx context.Context, provider, providerSubject string) (*coreentity.UserAuthIdentity, error)
	CreateAuthIdentity(ctx context.Context, identity coreentity.UserAuthIdentity) error
	UpdateAuthIdentityLastLogin(ctx context.Context, identityID string) error
	GetTenantProfile(ctx context.Context, tenantID string) (*coreentity.TenantProfile, error)
	GetTenantPreferredLanguage(ctx context.Context, tenantID string) (string, error)
	ExistsActiveUserByEmailAndTenant(ctx context.Context, email, tenantID string) (bool, error)
	ExistsActiveUserByEmailAndTenantExcludeUser(ctx context.Context, email, tenantID, excludeUserID string) (bool, error)
	GetRoleByName(ctx context.Context, roleName, tenantID string) (*coreentity.Role, error)
	InsertUser(ctx context.Context, user coreentity.User) (string, error)
	CreateUserActionToken(ctx context.Context, token coreentity.UserActionToken) error
	DeleteUserActionTokensByPurpose(ctx context.Context, userID, purpose string) error
	GetValidUserActionToken(ctx context.Context, tokenHash, purpose string) (*coreentity.UserActionToken, error)
	MarkUserActionTokenUsed(ctx context.Context, tokenID string) error
	MarkUserEmailVerified(ctx context.Context, userID string) error
	UpdateUserEmail(ctx context.Context, userID, tenantID, email string) error
	UpdateUserPassword(ctx context.Context, userID, tenantID, hashedPassword string) error
	GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error)
	GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error)
	CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error)
	UpdateUser(ctx context.Context, data coreentity.User) error
	DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error

	InitializeTenant(ctx context.Context, tenantID string, companyName string, preferredLanguage string) (string, error)
}
