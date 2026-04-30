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

func TestUserCore_CreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.User
		setup     func(repo *dbMocks.UserRepository)
		want      *coreentity.User
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.User{
				TenantID: "tenant-1",
				Email:    "newuser@example.com",
				Password: "password123",
				UserName: "New User",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "newuser@example.com", "tenant-1").
					Return(false, nil)
				repo.EXPECT().
					CreateUser(mock.Anything, mock.AnythingOfType("coreentity.User")).
					Return(&coreentity.User{
						ID:       "user-1",
						TenantID: "tenant-1",
						Email:    "newuser@example.com",
						UserName: "New User",
					}, nil)
			},
			want: &coreentity.User{
				ID:       "user-1",
				TenantID: "tenant-1",
				Email:    "newuser@example.com",
				UserName: "New User",
			},
		},
		{
			name: "duplicate email",
			input: coreentity.User{
				TenantID: "tenant-1",
				Email:    "existing@example.com",
				Password: "password123",
				UserName: "Existing User",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "existing@example.com", "tenant-1").
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - exists check",
			input: coreentity.User{
				TenantID: "tenant-1",
				Email:    "test@example.com",
				Password: "password123",
				UserName: "Test User",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "test@example.com", "tenant-1").
					Return(false, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - create user",
			input: coreentity.User{
				TenantID: "tenant-1",
				Email:    "test@example.com",
				Password: "password123",
				UserName: "Test User",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					ExistsActiveUserByEmailAndTenant(mock.Anything, "test@example.com", "tenant-1").
					Return(false, nil)
				repo.EXPECT().
					CreateUser(mock.Anything, mock.AnythingOfType("coreentity.User")).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, err := core.CreateUser(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.ID, got.ID)
			require.Equal(t, tt.want.TenantID, got.TenantID)
			require.Equal(t, tt.want.Email, got.Email)
			require.Equal(t, tt.want.UserName, got.UserName)
		})
	}
}
