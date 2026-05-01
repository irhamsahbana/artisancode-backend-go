package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_GetTenantProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			TenantID: "tenant-1",
		})
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			GetTenantProfile(mock.Anything, "tenant-1").
			Return(&coreentity.TenantProfile{
				ID:                  "tenant-1",
				Name:                "Acme Corp",
				Code:                "ACME",
				CanChangeTenantCode: true,
			}, nil)

		core := NewUserCore(Config{
			Repo: repo,
		})

		got, err := core.GetTenantProfile(ctx)

		require.NoError(t, err)
		require.Equal(t, "tenant-1", got.ID)
		require.Equal(t, "Acme Corp", got.Name)
		require.Equal(t, "ACME", got.Code)
		require.True(t, got.CanChangeTenantCode)
	})

	t.Run("rejects missing tenant context", func(t *testing.T) {
		core := NewUserCore(Config{})

		got, err := core.GetTenantProfile(context.Background())

		require.Nil(t, got)
		require.Error(t, err)

		customErr, ok := err.(*errmsg.CustomError)
		require.True(t, ok)
		require.Equal(t, errmsg.MessageInvalidCredentials, customErr.Msg)
	})

	t.Run("returns repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
			TenantID: "tenant-1",
		})
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			GetTenantProfile(mock.Anything, "tenant-1").
			Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))

		core := NewUserCore(Config{
			Repo: repo,
		})

		got, err := core.GetTenantProfile(ctx)

		require.Nil(t, got)
		require.Error(t, err)
	})
}
