package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.UserDeleteFilter
		setup     func(repo *dbMocks.UserRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.UserDeleteFilter{
				TenantID: "tenant-1",
				ID:       "user-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					DeleteUser(mock.Anything, coreentity.UserDeleteFilter{
						TenantID: "tenant-1",
						ID:       "user-1",
					}).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.UserDeleteFilter{
				TenantID: "tenant-1",
				ID:       "user-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					DeleteUser(mock.Anything, coreentity.UserDeleteFilter{
						TenantID: "tenant-1",
						ID:       "user-1",
					}).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)

			if tt.setup != nil {
				tt.setup(repo)
			}

			core := NewUserCore(Config{
				Repo: repo,
			})

			err := core.DeleteUser(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
