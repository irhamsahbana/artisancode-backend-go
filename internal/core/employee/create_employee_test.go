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

func TestEmployeeCore_CreateEmployee(t *testing.T) {
	ctx := context.Background()

	shiftID := "shift-1"
	joinDate := "2025-01-15"
	joinDateTZ := "Asia/Makassar"

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
				TenantID:         "tenant-1",
				EmployeeNo:       "EMP001",
				FullName:         "John Doe",
				ShiftID:          &shiftID,
				JoinDate:         &joinDate,
				JoinDateTimezone: &joinDateTZ,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "").
					Return(false, nil)
				repo.EXPECT().
					CreateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
					Return(&coreentity.Employee{
						ID:         "emp-1",
						TenantID:   "tenant-1",
						EmployeeNo: "EMP001",
						FullName:   "John Doe",
						ShiftID:    &shiftID,
						JoinDate:   &joinDate,
					}, nil)
			},
			want: &coreentity.Employee{
				ID:         "emp-1",
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				JoinDate:   &joinDate,
			},
		},
		{
			name: "invalid join date - timezone required",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				JoinDate:   &joinDate,
			},
			setup:     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name: "work shift is required",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
			},
			setup:     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name: "duplicate employee number",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - exists check",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "").
					Return(false, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - create",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "").
					Return(false, nil)
				repo.EXPECT().
					CreateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
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

			got, err := core.CreateEmployee(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// UserID should be nil (set to nil in core)
			require.Nil(t, got.UserID)
			require.Equal(t, tt.want.ID, got.ID)
			require.Equal(t, tt.want.TenantID, got.TenantID)
			require.Equal(t, tt.want.EmployeeNo, got.EmployeeNo)
			require.Equal(t, tt.want.FullName, got.FullName)
		})
	}
}
