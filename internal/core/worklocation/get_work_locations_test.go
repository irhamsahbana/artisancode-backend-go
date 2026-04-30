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

func TestWorkLocationCore_GetWorkLocations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    coreentity.WorkLocationListFilter
		setup     func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository)
		want      []coreentity.WorkLocation
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			filter: coreentity.WorkLocationListFilter{
				TenantID: "tenant-1",
				Q:        "main",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					GetWorkLocations(mock.Anything, coreentity.WorkLocationListFilter{
						TenantID: "tenant-1",
						Q:        "main",
						Page:     1,
						Paginate: 10,
					}).
					Return([]coreentity.WorkLocation{
						{ID: "wl-1", TenantID: "tenant-1", Name: "Main Office"},
					}, 1, nil)
			},
			want: []coreentity.WorkLocation{
				{ID: "wl-1", TenantID: "tenant-1", Name: "Main Office"},
			},
			wantCount: 1,
		},
		{
			name: "dependency error from repository",
			filter: coreentity.WorkLocationListFilter{
				TenantID: "tenant-1",
				Q:        "main",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					GetWorkLocations(mock.Anything, coreentity.WorkLocationListFilter{
						TenantID: "tenant-1",
						Q:        "main",
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
			repo := mocks.NewWorkLocationRepository(t)
			orgUnitRepo := mocks.NewOrgUnitRepository(t)

			if tt.setup != nil {
				tt.setup(repo, orgUnitRepo)
			}

			core := NewWorkLocationCore(Config{
				Repo:        repo,
				OrgUnitRepo: orgUnitRepo,
			})

			got, count, err := core.GetWorkLocations(ctx, tt.filter)

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
