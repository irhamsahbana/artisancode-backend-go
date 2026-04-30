package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRbacCore_GetRoles(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.RoleListFilter{
		TenantID: "tenant-1",
		Page:     1,
		Paginate: 10,
	}

	tests := []struct {
		name      string
		setup     func(repo *dbmocks.RbacRepository)
		want      []coreentity.Role
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			setup: func(repo *dbmocks.RbacRepository) {
				repo.EXPECT().
					GetRoles(mock.Anything, filter).
					Return([]coreentity.Role{{ID: "role-1", Name: "admin"}}, 1, nil)
			},
			want:      []coreentity.Role{{ID: "role-1", Name: "admin"}},
			wantCount: 1,
		},
		{
			name: "repository error",
			setup: func(repo *dbmocks.RbacRepository) {
				repo.EXPECT().
					GetRoles(mock.Anything, filter).
					Return(nil, 0, errors.New("repository failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewRbacRepository(t)
			tt.setup(repo)

			core := NewRbacCore(Config{Repo: repo})
			got, count, err := core.GetRoles(ctx, filter)

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
