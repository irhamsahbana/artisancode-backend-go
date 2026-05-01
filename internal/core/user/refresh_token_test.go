package core

import (
	"context"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_RefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success rotates refresh token and returns new access token", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		cache := tokencache.NewTokenCache(time.Hour, time.Minute)
		cache.SetRefreshToken(ctx, "old-refresh-token", tokencache.RefreshTokenData{
			UserID:   "user-1",
			TenantID: "tenant-1",
		}, time.Hour)

		repo.EXPECT().
			FindActiveUserByIDAndTenant(mock.Anything, "user-1", "tenant-1").
			Return(&coreentity.User{
				ID:          "user-1",
				TenantID:    "tenant-1",
				TenantName:  "Acme Corp",
				UserName:    "Test User",
				RoleNames:   []string{"owner"},
				CompanyID:   strPtr("company-1"),
				CompanyName: strPtr("Main Company"),
			}, nil)

		core := NewUserCore(Config{
			Repo:       repo,
			TokenCache: cache,
		})

		got, err := core.RefreshToken(ctx, coreentity.User{
			RefreshToken: "old-refresh-token",
		})

		require.NoError(t, err)
		require.NotEmpty(t, got.AccessToken)
		require.NotEmpty(t, got.RefreshToken)
		require.NotEqual(t, "old-refresh-token", got.RefreshToken)

		_, oldFound := cache.GetRefreshToken(ctx, "old-refresh-token")
		require.False(t, oldFound)

		newData, newFound := cache.GetRefreshToken(ctx, got.RefreshToken)
		require.True(t, newFound)
		require.Equal(t, "user-1", newData.UserID)
		require.Equal(t, "tenant-1", newData.TenantID)
	})

	t.Run("rejects missing refresh token", func(t *testing.T) {
		core := NewUserCore(Config{
			TokenCache: tokencache.NewTokenCache(time.Hour, time.Minute),
		})

		got, err := core.RefreshToken(ctx, coreentity.User{
			RefreshToken: "missing-refresh-token",
		})

		require.Nil(t, got)
		require.Error(t, err)

		customErr, ok := err.(*errmsg.CustomError)
		require.True(t, ok)
		require.Equal(t, errmsg.MessageInvalidOrExpiredRefreshToken, customErr.Msg)
	})

	t.Run("returns repository error", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		cache := tokencache.NewTokenCache(time.Hour, time.Minute)
		cache.SetRefreshToken(ctx, "valid-refresh-token", tokencache.RefreshTokenData{
			UserID:   "user-1",
			TenantID: "tenant-1",
		}, time.Hour)

		repo.EXPECT().
			FindActiveUserByIDAndTenant(mock.Anything, "user-1", "tenant-1").
			Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))

		core := NewUserCore(Config{
			Repo:       repo,
			TokenCache: cache,
		})

		got, err := core.RefreshToken(ctx, coreentity.User{
			RefreshToken: "valid-refresh-token",
		})

		require.Nil(t, got)
		require.Error(t, err)

		_, found := cache.GetRefreshToken(ctx, "valid-refresh-token")
		require.True(t, found)
	})
}
