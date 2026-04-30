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

func TestEmployeeCore_UpdateEmployee(t *testing.T) {
	ctx := context.Background()

	shiftID := "shift-1"
	userID := "user-1"

	tests := []struct {
		name      string
		input     coreentity.Employee
		setup     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository)
		wantError bool
	}{
		{
			name: "success - basic update without email or password change",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Updated",
				ShiftID:    &shiftID,
				Email:      "john@example.com",
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
						Email:      "john@example.com",
						UserID:     &userID,
					}, nil)
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "emp-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
					Return(nil)
			},
		},
		{
			name: "invalid join date - timezone required",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				JoinDate:   strPtr("2025-01-15"),
			},
			setup:     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name: "work shift is required",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
			},
			setup:     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {},
			wantError: true,
		},
		{
			name: "employee not found",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					GetEmployee(mock.Anything, coreentity.Employee{
						TenantID: "tenant-1",
						ID:       "emp-1",
					}).
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("employee not found"))
			},
			wantError: true,
		},
		{
			name: "duplicate employee number",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP002",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
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
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP002", "emp-1").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "email changed and user exists",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				Email:      "newemail@example.com",
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
						Email:      "oldemail@example.com",
						UserID:     &userID,
					}, nil)
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "emp-1").
					Return(false, nil)
				userRepo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "newemail@example.com", "tenant-1").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "success - update with email change and user has userID",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				Email:      "newemail@example.com",
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
						Email:      "oldemail@example.com",
						UserID:     &userID,
					}, nil)
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "emp-1").
					Return(false, nil)
				userRepo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "newemail@example.com", "tenant-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
					Return(nil)
				userRepo.EXPECT().
					UpdateUserEmail(mock.Anything, "user-1", "tenant-1", "newemail@example.com").
					Return(nil)
			},
		},
		{
			name: "success - update with password and user has userID",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
				Email:      "john@example.com",
				Password:   "newpassword123",
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
						Email:      "john@example.com",
						UserID:     &userID,
					}, nil)
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "emp-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
					Return(nil)
				userRepo.EXPECT().
					UpdateUserPassword(mock.Anything, "user-1", "tenant-1", mock.AnythingOfType("string")).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository - update",
			input: coreentity.Employee{
				TenantID:   "tenant-1",
				ID:         "emp-1",
				EmployeeNo: "EMP001",
				FullName:   "John Doe",
				ShiftID:    &shiftID,
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
				repo.EXPECT().
					ExistsByEmployeeNo(mock.Anything, "tenant-1", "EMP001", "emp-1").
					Return(false, nil)
				repo.EXPECT().
					UpdateEmployee(mock.Anything, mock.AnythingOfType("coreentity.Employee")).
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

			err := core.UpdateEmployee(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
