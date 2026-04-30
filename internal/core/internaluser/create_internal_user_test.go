package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateInternalUser(t *testing.T) {
	ctx := context.Background()
	userCtx := common.UserContext{Roles: []string{coreentity.InternalUserRoleSuperAdmin}}
	input := coreentity.InternalUser{
		UserCtx:  userCtx,
		FullName: " Admin User ",
		Email:    " ADMIN@EXAMPLE.COM ",
		Password: "secret",
		RoleCode: coreentity.InternalUserRoleOperator,
		Status:   coreentity.InternalUserStatusActive,
	}

	tests := []struct {
		name      string
		input     coreentity.InternalUser
		setup     func(repo *dbMocks.InternalUserRepository)
		wantError bool
	}{
		{
			name:  "normalizes and hashes password before create",
			input: input,
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					ExistsInternalUserByEmail(mock.Anything, "admin@example.com", "").
					Return(false, nil)
				repo.EXPECT().
					CreateInternalUser(
						mock.Anything,
						mock.MatchedBy(func(data coreentity.InternalUser) bool {
							return data.FullName == "Admin User" &&
								data.Email == "admin@example.com" &&
								data.Password != "secret" &&
								data.RoleCode == coreentity.InternalUserRoleOperator
						}),
					).
					Return(&coreentity.InternalUser{ID: "user-1"}, nil)
			},
		},
		{
			name: "rejects unauthorized user before repository call",
			input: coreentity.InternalUser{
				UserCtx:  common.UserContext{Roles: []string{coreentity.InternalUserRoleOperator}},
				Email:    "admin@example.com",
				Password: "secret",
				RoleCode: coreentity.InternalUserRoleOperator,
				Status:   coreentity.InternalUserStatusActive,
			},
			wantError: true,
		},
		{
			name: "rejects missing password",
			input: coreentity.InternalUser{
				UserCtx:  userCtx,
				Email:    "admin@example.com",
				RoleCode: coreentity.InternalUserRoleOperator,
				Status:   coreentity.InternalUserStatusActive,
			},
			wantError: true,
		},
		{
			name:  "rejects duplicate email",
			input: input,
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					ExistsInternalUserByEmail(mock.Anything, "admin@example.com", "").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name:  "returns repository error",
			input: input,
			setup: func(repo *dbMocks.InternalUserRepository) {
				repo.EXPECT().
					ExistsInternalUserByEmail(mock.Anything, "admin@example.com", "").
					Return(false, errors.New("repo failed"))
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
			core := NewInternalUserCore(Config{Repo: repo})

			got, err := core.CreateInternalUser(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
