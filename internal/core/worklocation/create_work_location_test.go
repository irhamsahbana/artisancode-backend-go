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

func TestWorkLocationCore_CreateWorkLocation(t *testing.T) {
	ctx := context.Background()
	orgUnitID := "org-1"

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
				Name:     "Main Office",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Main Office", (*string)(nil)).
					Return(false, nil)
				repo.EXPECT().
					CreateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
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
			name: "success with org unit",
			input: coreentity.WorkLocation{
				TenantID:  "tenant-1",
				Name:      "Branch Office",
				OrgUnitID: &orgUnitID,
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Branch Office", (*string)(nil)).
					Return(false, nil)
				orgUnitRepo.EXPECT().
					GetOrgUnit(mock.Anything, coreentity.OrgUnit{
						TenantID: "tenant-1",
						ID:       "org-1",
					}).
					Return(&coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1"}, nil)
				repo.EXPECT().
					CreateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
					Return(&coreentity.WorkLocation{
						ID:        "wl-2",
						TenantID:  "tenant-1",
						Name:      "Branch Office",
						OrgUnitID: &orgUnitID,
					}, nil)
			},
			want: &coreentity.WorkLocation{
				ID:        "wl-2",
				TenantID:  "tenant-1",
				Name:      "Branch Office",
				OrgUnitID: &orgUnitID,
			},
		},
		{
			name: "duplicate name error",
			input: coreentity.WorkLocation{
				TenantID: "tenant-1",
				Name:     "Main Office",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Main Office", (*string)(nil)).
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "org unit not found",
			input: coreentity.WorkLocation{
				TenantID:  "tenant-1",
				Name:      "Branch Office",
				OrgUnitID: &orgUnitID,
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Branch Office", (*string)(nil)).
					Return(false, nil)
				orgUnitRepo.EXPECT().
					GetOrgUnit(mock.Anything, coreentity.OrgUnit{
						TenantID: "tenant-1",
						ID:       "org-1",
					}).
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("org unit not found"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - exists check",
			input: coreentity.WorkLocation{
				TenantID: "tenant-1",
				Name:     "Main Office",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Main Office", (*string)(nil)).
					Return(false, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - create",
			input: coreentity.WorkLocation{
				TenantID: "tenant-1",
				Name:     "Main Office",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Main Office", (*string)(nil)).
					Return(false, nil)
				repo.EXPECT().
					CreateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
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

			got, err := core.CreateWorkLocation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
