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

func TestUserCore_GetUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    coreentity.UserListFilter
		setup     func(repo *dbMocks.UserRepository)
		want      []coreentity.User
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			filter: coreentity.UserListFilter{
				TenantID: "tenant-1",
				Q:        "test",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					GetUsers(mock.Anything, coreentity.UserListFilter{
						TenantID: "tenant-1",
						Q:        "test",
						Page:     1,
						Paginate: 10,
					}).
					Return([]coreentity.User{
						{
							ID:       "user-1",
							TenantID: "tenant-1",
							Email:    "test@example.com",
							UserName: "Test User",
						},
					}, 1, nil)
			},
			want: []coreentity.User{
				{
					ID:       "user-1",
					TenantID: "tenant-1",
					Email:    "test@example.com",
					UserName: "Test User",
				},
			},
			wantCount: 1,
		},
		{
			name: "dependency error from repository",
			filter: coreentity.UserListFilter{
				TenantID: "tenant-1",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					GetUsers(mock.Anything, mock.AnythingOfType("coreentity.UserListFilter")).
					Return(nil, 0, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, count, err := core.GetUsers(ctx, tt.filter)

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
