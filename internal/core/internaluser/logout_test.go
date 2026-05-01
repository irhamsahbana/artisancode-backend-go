package core

import (
	"context"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"

	"github.com/stretchr/testify/require"
)

func TestInternalUserCore_Logout(t *testing.T) {
	t.Run("deletes current refresh token", func(t *testing.T) {
		ctx := context.Background()
		cache := tokencache.NewTokenCache(time.Hour, time.Minute)
		cache.SetRefreshToken(ctx, "refresh-token", tokencache.RefreshTokenData{
			UserID: "user-1",
		}, time.Hour)

		core := NewInternalUserCore(Config{TokenCache: cache})

		err := core.Logout(ctx, coreentity.InternalUser{
			RefreshToken: "refresh-token",
		})

		require.NoError(t, err)

		_, found := cache.GetRefreshToken(ctx, "refresh-token")
		require.False(t, found)
	})

	t.Run("accepts missing refresh token", func(t *testing.T) {
		ctx := context.Background()
		core := NewInternalUserCore(Config{
			TokenCache: tokencache.NewTokenCache(time.Hour, time.Minute),
		})

		err := core.Logout(ctx, coreentity.InternalUser{})

		require.NoError(t, err)
	})
}
