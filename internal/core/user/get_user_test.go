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

func TestUserCore_GetUser(t *testing.T) {
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
				ID:       "user-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					GetUser(mock.Anything, coreentity.User{
						TenantID: "tenant-1",
						ID:       "user-1",
					}).
					Return(&coreentity.User{
						ID:       "user-1",
						TenantID: "tenant-1",
						Email:    "test@example.com",
						UserName: "Test User",
					}, nil)
			},
			want: &coreentity.User{
				ID:       "user-1",
				TenantID: "tenant-1",
				Email:    "test@example.com",
				UserName: "Test User",
			},
		},
		{
			name: "user not found",
			input: coreentity.User{
				TenantID: "tenant-1",
				ID:       "nonexistent",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					GetUser(mock.Anything, coreentity.User{
						TenantID: "tenant-1",
						ID:       "nonexistent",
					}).
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("user not found"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository",
			input: coreentity.User{
				TenantID: "tenant-1",
				ID:       "user-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					GetUser(mock.Anything, coreentity.User{
						TenantID: "tenant-1",
						ID:       "user-1",
					}).
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

			got, err := core.GetUser(ctx, tt.input)

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
