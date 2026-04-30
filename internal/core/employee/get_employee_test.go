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

func TestEmployeeCore_GetEmployee(t *testing.T) {
	ctx := context.Background()

	shiftID := "shift-1"

	tests := []struct {
		name      string
		input     coreentity.Employee
		setup     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository)
		want      *coreentity.Employee
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.Employee{
				TenantID: "tenant-1",
				ID:       "emp-1",
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					GetEmployee(mock.Anything, coreentity.Employee{
						TenantID: "tenant-1",
						ID:       "emp-1",
					}).
					Return(&coreentity.Employee{
						ID:         "emp-1",
						TenantID:   "tenant-1",
						EmployeeNo: "EMP001",
						FullName:   "John Doe",
						ShiftID:    &shiftID,
					}, nil)
			},
			want: &coreentity.Employee{
				ID:         "emp-1",
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.Employee{
				TenantID: "tenant-1",
				ID:       "emp-1",
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					GetEmployee(mock.Anything, coreentity.Employee{
						TenantID: "tenant-1",
						ID:       "emp-1",
					}).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, err := core.GetEmployee(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
