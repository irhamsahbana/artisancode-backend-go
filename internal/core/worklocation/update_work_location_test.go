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

func TestWorkLocationCore_UpdateWorkLocation(t *testing.T) {
	ctx := context.Background()
	orgUnitID := "org-1"
	wlID := "wl-1"

	tests := []struct {
		name      string
		input     coreentity.WorkLocation
		setup     func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkLocation{
				ID:       "wl-1",
				TenantID: "tenant-1",
				Name:     "Updated Office",
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
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Updated Office", &wlID).
					Return(false, nil)
				repo.EXPECT().
					UpdateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
					Return(nil)
			},
		},
		{
			name: "success with org unit",
			input: coreentity.WorkLocation{
				ID:        "wl-1",
				TenantID:  "tenant-1",
				Name:      "Updated Office",
				OrgUnitID: &orgUnitID,
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
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Updated Office", &wlID).
					Return(false, nil)
				orgUnitRepo.EXPECT().
					GetOrgUnit(mock.Anything, coreentity.OrgUnit{
						TenantID: "tenant-1",
						ID:       "org-1",
					}).
					Return(&coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1"}, nil)
				repo.EXPECT().
					UpdateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
					Return(nil)
			},
		},
		{
			name: "work location not found",
			input: coreentity.WorkLocation{
				ID:       "wl-1",
				TenantID: "tenant-1",
				Name:     "Updated Office",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					GetWorkLocation(mock.Anything, coreentity.WorkLocation{
						TenantID: "tenant-1",
						ID:       "wl-1",
					}).
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("not found"))
			},
			wantError: true,
		},
		{
			name: "duplicate name error",
			input: coreentity.WorkLocation{
				ID:       "wl-1",
				TenantID: "tenant-1",
				Name:     "Updated Office",
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
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Updated Office", &wlID).
					Return(true, nil)
			},
			wantError: true,
		},
		{
			name: "org unit not found",
			input: coreentity.WorkLocation{
				ID:        "wl-1",
				TenantID:  "tenant-1",
				Name:      "Updated Office",
				OrgUnitID: &orgUnitID,
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
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Updated Office", &wlID).
					Return(false, nil)
				orgUnitRepo.EXPECT().
					GetOrgUnit(mock.Anything, coreentity.OrgUnit{
						TenantID: "tenant-1",
						ID:       "org-1",
					}).
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("not found"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository - update",
			input: coreentity.WorkLocation{
				ID:       "wl-1",
				TenantID: "tenant-1",
				Name:     "Updated Office",
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
				repo.EXPECT().
					ExistsWorkLocationByName(mock.Anything, "tenant-1", "Updated Office", &wlID).
					Return(false, nil)
				repo.EXPECT().
					UpdateWorkLocation(mock.Anything, mock.AnythingOfType("coreentity.WorkLocation")).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			err := core.UpdateWorkLocation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
