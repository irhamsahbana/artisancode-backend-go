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

func TestUpdateInternalUser(t *testing.T) {
	ctx := context.Background()
	input := coreentity.InternalUser{
		UserCtx:  common.UserContext{Roles: []string{coreentity.InternalUserRoleSuperAdmin}},
		ID:       "user-1",
		FullName: " Admin User ",
		Email:    " ADMIN@EXAMPLE.COM ",
		Password: "new-secret",
		RoleCode: coreentity.InternalUserRoleOperator,
		Status:   coreentity.InternalUserStatusActive,
	}
	repo := dbMocks.NewInternalUserRepository(t)
	repo.EXPECT().
		ExistsInternalUserByEmail(mock.Anything, "admin@example.com", "user-1").
		Return(false, nil)
	repo.EXPECT().
		UpdateInternalUser(mock.Anything, mock.MatchedBy(func(data coreentity.InternalUser) bool {
			return data.ID == "user-1" &&
				data.FullName == "Admin User" &&
				data.Email == "admin@example.com" &&
				data.Password != "new-secret"
		})).
		Return(nil)
	core := NewInternalUserCore(Config{Repo: repo})

	err := core.UpdateInternalUser(ctx, input)

	require.NoError(t, err)
}
