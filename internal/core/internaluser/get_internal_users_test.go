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

func TestGetInternalUsers(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
		Roles: []string{coreentity.InternalUserRoleSuperAdmin},
	})
	filter := coreentity.InternalUserListFilter{RoleCode: coreentity.InternalUserRoleOperator}
	repo := dbMocks.NewInternalUserRepository(t)
	repo.EXPECT().
		GetInternalUsers(mock.Anything, filter).
		Return([]coreentity.InternalUser{{ID: "user-1"}}, 1, nil)
	core := NewInternalUserCore(Config{Repo: repo})

	got, total, err := core.GetInternalUsers(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, []coreentity.InternalUser{{ID: "user-1"}}, got)
	require.Equal(t, 1, total)
}
