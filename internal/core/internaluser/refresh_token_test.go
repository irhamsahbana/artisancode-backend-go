package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(cache *tokencache.TokenCache, repo *dbMocks.InternalUserRepository)
		wantError bool
	}{
		{
			name: "rotates valid refresh token",
			setup: func(cache *tokencache.TokenCache, repo *dbMocks.InternalUserRepository) {
				cache.SetRefreshToken("refresh-1", tokencache.RefreshTokenData{UserID: "user-1"}, time.Hour)
				repo.EXPECT().
					GetInternalUser(mock.Anything, coreentity.InternalUserFilter{ID: "user-1"}).
					Return(&coreentity.InternalUser{
						ID:       "user-1",
						FullName: "Admin",
						RoleCode: coreentity.InternalUserRoleOperator,
						Status:   coreentity.InternalUserStatusActive,
					}, nil)
			},
		},
		{
			name:      "rejects missing refresh token",
			wantError: true,
		},
		{
			name: "returns repository error",
			setup: func(cache *tokencache.TokenCache, repo *dbMocks.InternalUserRepository) {
				cache.SetRefreshToken("refresh-1", tokencache.RefreshTokenData{UserID: "user-1"}, time.Hour)
				repo.EXPECT().
					GetInternalUser(mock.Anything, coreentity.InternalUserFilter{ID: "user-1"}).
					Return(nil, errors.New("repo failed"))
			},
			wantError: true,
		},
		{
			name: "rejects inactive internal user",
			setup: func(cache *tokencache.TokenCache, repo *dbMocks.InternalUserRepository) {
				cache.SetRefreshToken("refresh-1", tokencache.RefreshTokenData{UserID: "user-1"}, time.Hour)
				repo.EXPECT().
					GetInternalUser(mock.Anything, coreentity.InternalUserFilter{ID: "user-1"}).
					Return(&coreentity.InternalUser{
						ID:     "user-1",
						Status: coreentity.InternalUserStatusInactive,
					}, nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalUserRepository(t)
			cache := tokencache.NewTokenCache(time.Hour, time.Hour)
			if tt.setup != nil {
				tt.setup(cache, repo)
			}
			core := NewInternalUserCore(Config{Repo: repo, TokenCache: cache})

			got, err := core.RefreshToken(ctx, coreentity.InternalUser{RefreshToken: "refresh-1"})

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.AccessToken)
			require.NotEmpty(t, got.RefreshToken)
			require.NotEqual(t, "refresh-1", got.RefreshToken)
			_, oldFound := cache.GetRefreshToken("refresh-1")
			_, newFound := cache.GetRefreshToken(got.RefreshToken)
			require.False(t, oldFound)
			require.True(t, newFound)
		})
	}
}
