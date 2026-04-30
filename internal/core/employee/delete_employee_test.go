package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmployeeCore_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.EmployeeDeleteFilter
		setup     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.EmployeeDeleteFilter{
				TenantID: "tenant-1",
				ID:       "emp-1",
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					DeleteEmployee(mock.Anything, coreentity.EmployeeDeleteFilter{
						TenantID: "tenant-1",
						ID:       "emp-1",
					}).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.EmployeeDeleteFilter{
				TenantID: "tenant-1",
				ID:       "emp-1",
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					DeleteEmployee(mock.Anything, coreentity.EmployeeDeleteFilter{
						TenantID: "tenant-1",
						ID:       "emp-1",
					}).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewEmployeeRepository(t)
			userRepo := mocks.NewUserRepository(t)

			if tt.setup != nil {
				tt.setup(repo, userRepo)
			}

			core := NewEmployeeCore(Config{
				Repo:     repo,
				UserRepo: userRepo,
			})

			err := core.DeleteEmployee(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
