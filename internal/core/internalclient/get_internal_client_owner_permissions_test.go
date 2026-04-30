package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInternalClientOwnerPermissions(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
		Roles: []string{coreentity.InternalUserRoleSuperAdmin},
	})
	repo := dbMocks.NewInternalClientRepository(t)
	repo.EXPECT().
		GetInternalClientOwnerPermissions(mock.Anything, "client-1").
		Return(&coreentity.InternalClientOwnerPermissions{ClientID: "client-1"}, nil)
	core := NewInternalClientCore(Config{Repo: repo})

	got, err := core.GetInternalClientOwnerPermissions(ctx, "client-1")

	require.NoError(t, err)
	require.Equal(t, "client-1", got.ClientID)
}
