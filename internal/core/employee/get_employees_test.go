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

func TestEmployeeCore_GetEmployees(t *testing.T) {
	ctx := context.Background()

	shiftID := "shift-1"

	tests := []struct {
		name      string
		filter    coreentity.EmployeeListFilter
		setup     func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository)
		want      []coreentity.Employee
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			filter: coreentity.EmployeeListFilter{
				TenantID: "tenant-1",
				Q:        "john",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					GetEmployees(mock.Anything, coreentity.EmployeeListFilter{
						TenantID: "tenant-1",
						Q:        "john",
						Page:     1,
						Paginate: 10,
					}).
					Return([]coreentity.Employee{
						{
							ID:         "emp-1",
							TenantID:   "tenant-1",
							EmployeeNo: "EMP001",
							FullName:   "John Doe",
							ShiftID:    &shiftID,
						},
					}, 1, nil)
			},
			want: []coreentity.Employee{
				{
					ID:         "emp-1",
					TenantID:   "tenant-1",
					EmployeeNo: "EMP001",
					FullName:   "John Doe",
					ShiftID:    &shiftID,
				},
			},
			wantCount: 1,
		},
		{
			name: "dependency error from repository",
			filter: coreentity.EmployeeListFilter{
				TenantID: "tenant-1",
				Q:        "john",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.EmployeeRepository, userRepo *mocks.UserRepository) {
				repo.EXPECT().
					GetEmployees(mock.Anything, coreentity.EmployeeListFilter{
						TenantID: "tenant-1",
						Q:        "john",
						Page:     1,
						Paginate: 10,
					}).
					Return(nil, 0, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, count, err := core.GetEmployees(ctx, tt.filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantCount, count)
		})
	}
}
