package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type UserCore interface {
	Login(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	GoogleRegister(ctx context.Context, input coreentity.GoogleRegisterInput) (*coreentity.GoogleRegisterResult, error)
	GoogleLogin(ctx context.Context, input coreentity.GoogleLoginInput) (*coreentity.AuthTokens, error)
	RefreshToken(ctx context.Context, user coreentity.User) (*coreentity.AuthTokens, error)
	RegisterOwner(ctx context.Context, user coreentity.User, tenant coreentity.Tenant) (*coreentity.RegisterResult, error)
	GetTenantProfile(ctx context.Context) (*coreentity.TenantProfile, error)
	VerifyEmail(ctx context.Context, token coreentity.UserActionToken) error
	ResendVerificationEmail(ctx context.Context, user coreentity.User) error
	ForgotPassword(ctx context.Context, user coreentity.User) error
	ResetPassword(ctx context.Context, token coreentity.UserActionToken, user coreentity.User) error
	GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error)
	GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error)
	CreateUser(ctx context.Context, data coreentity.User) (*coreentity.User, error)
	UpdateUser(ctx context.Context, data coreentity.User) error
	DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error
}
