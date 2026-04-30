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

func TestUpdateInternalClientOwnerPermissions(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
		Roles: []string{coreentity.InternalUserRoleSuperAdmin},
	})
	input := coreentity.InternalClientOwnerPermissionUpdate{
		ClientID:      "client-1",
		PermissionIDs: []string{" permission-1 ", "", "permission-1", "permission-2"},
	}

	tests := []struct {
		name      string
		ctx       context.Context
		setup     func(repo *dbMocks.InternalClientRepository)
		wantError bool
	}{
		{
			name: "normalizes unique permission ids",
			ctx:  ctx,
			setup: func(repo *dbMocks.InternalClientRepository) {
				repo.EXPECT().
					UpdateInternalClientOwnerPermissions(
						mock.Anything,
						coreentity.InternalClientOwnerPermissionUpdate{
							ClientID:      "client-1",
							PermissionIDs: []string{"permission-1", "permission-2"},
						},
					).
					Return(nil)
			},
		},
		{
			name:      "rejects unauthorized user before repository call",
			ctx:       context.Background(),
			wantError: true,
		},
		{
			name: "returns repository error",
			ctx:  ctx,
			setup: func(repo *dbMocks.InternalClientRepository) {
				repo.EXPECT().
					UpdateInternalClientOwnerPermissions(mock.Anything, mock.Anything).
					Return(errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalClientRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			core := NewInternalClientCore(Config{Repo: repo})

			err := core.UpdateInternalClientOwnerPermissions(tt.ctx, input)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
