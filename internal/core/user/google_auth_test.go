package core

import (
	"context"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_GoogleRegister(t *testing.T) {
	ctx := context.Background()
	identity := &coreentity.GoogleIdentity{
		Subject:       "google-subject-1",
		Email:         "owner@example.com",
		EmailVerified: true,
		DisplayName:   "Owner Example",
	}

	t.Run("success creates tenant user identity and returns tokens", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		validator := integrationMocks.NewGoogleIDTokenValidator(t)

		validator.EXPECT().
			Validate(mock.Anything, "id-token", "").
			Return(identity, nil)
		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().ExistsTenantByCode(mock.Anything, "GOOG").Return(false, nil)
		repo.EXPECT().
			FindAuthIdentityByProviderSubject(mock.Anything, coreentity.AuthProviderGoogle, identity.Subject).
			Return(nil, nil)
		repo.EXPECT().FindActiveUsersByEmail(mock.Anything, identity.Email).Return(nil, nil)
		repo.EXPECT().
			InsertTenant(mock.Anything, coreentity.Tenant{Name: "PT Contoh", Code: "GOOG"}).
			Return("tenant-1", nil)
		repo.EXPECT().
			InitializeTenant(mock.Anything, "tenant-1", "PT Contoh", "id").
			Return("company-1", nil)
		repo.EXPECT().
			GetRoleByName(mock.Anything, "owner", "tenant-1").
			Return(&coreentity.Role{ID: "role-owner", Name: "owner"}, nil)
		repo.EXPECT().
			InsertUser(mock.Anything, mock.MatchedBy(func(user coreentity.User) bool {
				return user.Email == identity.Email &&
					user.TenantID == "tenant-1" &&
					user.EmailVerifiedAt != nil
			})).
			Return("user-1", nil)
		repo.EXPECT().
			CreateAuthIdentity(mock.Anything, mock.MatchedBy(func(auth coreentity.UserAuthIdentity) bool {
				return auth.UserID == "user-1" &&
					auth.TenantID == "tenant-1" &&
					auth.Provider == coreentity.AuthProviderGoogle &&
					auth.ProviderSubject == identity.Subject
			})).
			Return(nil)

		core := NewUserCore(Config{
			Repo:                 repo,
			Tx:                   tx,
			TokenCache:           tokencache.NewTokenCache(time.Hour, time.Minute),
			GoogleTokenValidator: validator,
		})

		got, err := core.GoogleRegister(ctx, coreentity.GoogleRegisterInput{
			IDToken:            "id-token",
			TenantName:         "PT Contoh",
			TenantCode:         "goog",
			ConfirmTenantSetup: true,
			PreferredLanguage:  "id",
		})

		require.NoError(t, err)
		require.NotEmpty(t, got.AccessToken)
		require.NotEmpty(t, got.RefreshToken)
		require.Equal(t, "GOOG", got.TenantCode)
	})

	t.Run("rejects missing confirmation before token validation", func(t *testing.T) {
		core := NewUserCore(Config{})

		got, err := core.GoogleRegister(ctx, coreentity.GoogleRegisterInput{
			IDToken:    "id-token",
			TenantCode: "T2K2",
		})

		require.Nil(t, got)
		require.Error(t, err)
	})

	t.Run("rejects reserved tenant code", func(t *testing.T) {
		core := NewUserCore(Config{})

		got, err := core.GoogleRegister(ctx, coreentity.GoogleRegisterInput{
			IDToken:            "id-token",
			TenantCode:         "ADMIN",
			ConfirmTenantSetup: true,
		})

		require.Nil(t, got)
		require.Error(t, err)
	})

	t.Run("success consumes registration token without validating id token again", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		cache := tokencache.NewTokenCache(time.Hour, time.Minute)
		cache.SetGoogleRegistration(ctx, "reg-token", tokencache.GoogleRegistrationData{
			Subject:       identity.Subject,
			Email:         identity.Email,
			EmailVerified: identity.EmailVerified,
			DisplayName:   identity.DisplayName,
		}, time.Hour)

		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().ExistsTenantByCode(mock.Anything, "GOOG").Return(false, nil)
		repo.EXPECT().
			FindAuthIdentityByProviderSubject(mock.Anything, coreentity.AuthProviderGoogle, identity.Subject).
			Return(nil, nil)
		repo.EXPECT().FindActiveUsersByEmail(mock.Anything, identity.Email).Return(nil, nil)
		repo.EXPECT().
			InsertTenant(mock.Anything, coreentity.Tenant{Name: "PT Contoh", Code: "GOOG"}).
			Return("tenant-1", nil)
		repo.EXPECT().
			InitializeTenant(mock.Anything, "tenant-1", "PT Contoh", "id").
			Return("company-1", nil)
		repo.EXPECT().
			GetRoleByName(mock.Anything, "owner", "tenant-1").
			Return(&coreentity.Role{ID: "role-owner", Name: "owner"}, nil)
		repo.EXPECT().
			InsertUser(mock.Anything, mock.Anything).
			Return("user-1", nil)
		repo.EXPECT().
			CreateAuthIdentity(mock.Anything, mock.Anything).
			Return(nil)

		core := NewUserCore(Config{
			Repo:       repo,
			Tx:         tx,
			TokenCache: cache,
		})

		got, err := core.GoogleRegister(ctx, coreentity.GoogleRegisterInput{
			RegistrationToken:  "reg-token",
			TenantName:         "PT Contoh",
			TenantCode:         "goog",
			ConfirmTenantSetup: true,
			PreferredLanguage:  "id",
		})

		require.NoError(t, err)
		require.NotEmpty(t, got.AccessToken)
		require.False(t, hasGoogleRegistration(ctx, cache, "reg-token"))
	})

	t.Run("rejects invalid registration token", func(t *testing.T) {
		core := NewUserCore(Config{
			TokenCache: tokencache.NewTokenCache(time.Hour, time.Minute),
		})

		got, err := core.GoogleRegister(ctx, coreentity.GoogleRegisterInput{
			RegistrationToken:  "missing",
			TenantName:         "PT Contoh",
			TenantCode:         "T2K2",
			ConfirmTenantSetup: true,
		})

		require.Nil(t, got)
		require.Error(t, err)
	})
}

func TestUserCore_GoogleRegisterInit(t *testing.T) {
	ctx := context.Background()
	repo := dbMocks.NewUserRepository(t)
	tx := dbMocks.NewTransactor(t)
	validator := integrationMocks.NewGoogleIDTokenValidator(t)
	cache := tokencache.NewTokenCache(time.Hour, time.Minute)

	validator.EXPECT().
		Validate(mock.Anything, "id-token", "").
		Return(&coreentity.GoogleIdentity{
			Subject:       "google-subject-1",
			Email:         "owner@example.com",
			EmailVerified: true,
			DisplayName:   "Owner Example",
		}, nil)

	core := NewUserCore(Config{
		Repo:                 repo,
		Tx:                   tx,
		TokenCache:           cache,
		GoogleTokenValidator: validator,
	})

	got, err := core.GoogleRegisterInit(ctx, coreentity.GoogleRegisterInitInput{
		IDToken: "id-token",
	})

	require.NoError(t, err)
	require.NotEmpty(t, got.RegistrationToken)
	require.Equal(t, "owner@example.com", got.Email)
	require.True(t, hasGoogleRegistration(ctx, cache, got.RegistrationToken))
}

func TestUserCore_GoogleLogin(t *testing.T) {
	ctx := context.Background()
	identity := &coreentity.GoogleIdentity{
		Subject:       "google-subject-1",
		Email:         "user@example.com",
		EmailVerified: true,
		DisplayName:   "User Example",
	}

	t.Run("auto-links when email has exactly one active user", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		validator := integrationMocks.NewGoogleIDTokenValidator(t)

		validator.EXPECT().
			Validate(mock.Anything, "id-token", "").
			Return(identity, nil)
		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().
			FindAuthIdentityByProviderSubject(mock.Anything, coreentity.AuthProviderGoogle, identity.Subject).
			Return(nil, nil)
		repo.EXPECT().
			FindActiveUsersByEmail(mock.Anything, identity.Email).
			Return([]coreentity.User{{
				ID:         "user-1",
				TenantID:   "tenant-1",
				TenantName: "PT Contoh",
				UserName:   "user",
				Email:      identity.Email,
				RoleNames:  []string{"employee"},
			}}, nil)
		repo.EXPECT().
			CreateAuthIdentity(mock.Anything, mock.MatchedBy(func(auth coreentity.UserAuthIdentity) bool {
				return auth.UserID == "user-1" &&
					auth.ProviderSubject == identity.Subject
			})).
			Return(nil)

		core := NewUserCore(Config{
			Repo:                 repo,
			Tx:                   tx,
			TokenCache:           tokencache.NewTokenCache(time.Hour, time.Minute),
			GoogleTokenValidator: validator,
		})

		got, err := core.GoogleLogin(ctx, coreentity.GoogleLoginInput{IDToken: "id-token"})

		require.NoError(t, err)
		require.NotEmpty(t, got.AccessToken)
		require.NotEmpty(t, got.RefreshToken)
	})

	t.Run("rejects ambiguous active email", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		validator := integrationMocks.NewGoogleIDTokenValidator(t)

		validator.EXPECT().
			Validate(mock.Anything, "id-token", "").
			Return(identity, nil)
		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().
			FindAuthIdentityByProviderSubject(mock.Anything, coreentity.AuthProviderGoogle, identity.Subject).
			Return(nil, nil)
		repo.EXPECT().
			FindActiveUsersByEmail(mock.Anything, identity.Email).
			Return([]coreentity.User{
				{ID: "user-1", TenantID: "tenant-1"},
				{ID: "user-2", TenantID: "tenant-2"},
			}, nil)

		core := NewUserCore(Config{
			Repo:                 repo,
			Tx:                   tx,
			TokenCache:           tokencache.NewTokenCache(time.Hour, time.Minute),
			GoogleTokenValidator: validator,
		})

		got, err := core.GoogleLogin(ctx, coreentity.GoogleLoginInput{IDToken: "id-token"})

		require.Nil(t, got)
		require.Error(t, err)
	})
}

func hasGoogleRegistration(ctx context.Context, cache *tokencache.TokenCache, token string) bool {
	_, found := cache.GetGoogleRegistration(ctx, token)
	return found
}
