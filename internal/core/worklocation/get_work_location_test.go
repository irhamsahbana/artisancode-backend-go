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

func TestWorkLocationCore_GetWorkLocation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.WorkLocation
		setup     func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository)
		want      *coreentity.WorkLocation
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkLocation{
				TenantID: "tenant-1",
				ID:       "wl-1",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					GetWorkLocation(mock.Anything, coreentity.WorkLocation{
						TenantID: "tenant-1",
						ID:       "wl-1",
					}).
					Return(&coreentity.WorkLocation{
						ID:       "wl-1",
						TenantID: "tenant-1",
						Name:     "Main Office",
					}, nil)
			},
			want: &coreentity.WorkLocation{
				ID:       "wl-1",
				TenantID: "tenant-1",
				Name:     "Main Office",
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkLocation{
				TenantID: "tenant-1",
				ID:       "wl-1",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					GetWorkLocation(mock.Anything, coreentity.WorkLocation{
						TenantID: "tenant-1",
						ID:       "wl-1",
					}).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, err := core.GetWorkLocation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
