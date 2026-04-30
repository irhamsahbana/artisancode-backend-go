package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	"codebase-app/internal/integration/tokencache"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "internaluser-core-test"
	infraConfig.Envs.Guard.JwtPrivateKey = "test-secret"
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	hashedPassword, err := hashInternalPassword("secret")
	require.NoError(t, err)

	tests := []struct {
		name      string
		input     coreentity.InternalUser
		setup     func(repo *dbMocks.InternalUserRepository)
		wantError bool
	}{
		{
			name: coreentity.InternalUserStatusActive + " user gets tokens and last login update",
			input: coreentity.InternalUser{
				Email:    " ADMIN@EXAMPLE.COM ",
				Password: "secret",
			},
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					FindActiveInternalUserByEmail(mock.Anything, "admin@example.com").
					Return(&coreentity.InternalUser{
						ID:       "user-1",
						FullName: "Admin",
						Email:    "admin@example.com",
						Password: hashedPassword,
						RoleCode: coreentity.InternalUserRoleSuperAdmin,
						Status:   coreentity.InternalUserStatusActive,
					}, nil)
				repo.EXPECT().
					UpdateInternalUserLastLogin(mock.Anything, "user-1").
					Return(nil)
			},
		},
		{
			name: "rejects invalid password",
			input: coreentity.InternalUser{
				Email:    "admin@example.com",
				Password: "wrong",
			},
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					FindActiveInternalUserByEmail(mock.Anything, "admin@example.com").
					Return(&coreentity.InternalUser{
						ID:       "user-1",
						Password: hashedPassword,
						Status:   coreentity.InternalUserStatusActive,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "rejects inactive user",
			input: coreentity.InternalUser{
				Email:    "admin@example.com",
				Password: "secret",
			},
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					FindActiveInternalUserByEmail(mock.Anything, "admin@example.com").
					Return(&coreentity.InternalUser{
						ID:       "user-1",
						Password: hashedPassword,
						Status:   coreentity.InternalUserStatusInactive,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "returns repository lookup error",
			input: coreentity.InternalUser{
				Email:    "admin@example.com",
				Password: "secret",
			},
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					FindActiveInternalUserByEmail(mock.Anything, "admin@example.com").
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalUserRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			cache := tokencache.NewTokenCache(time.Hour, time.Hour)
			core := NewInternalUserCore(Config{Repo: repo, TokenCache: cache})

			got, err := core.Login(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.AccessToken)
			require.NotEmpty(t, got.RefreshToken)
			_, found := cache.GetRefreshToken(got.RefreshToken)
			require.True(t, found)
		})
	}
}
