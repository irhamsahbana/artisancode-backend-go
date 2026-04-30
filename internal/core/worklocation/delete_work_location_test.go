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

func TestWorkLocationCore_DeleteWorkLocation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.WorkLocationDeleteFilter
		setup     func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkLocationDeleteFilter{
				TenantID: "tenant-1",
				ID:       "wl-1",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					DeleteWorkLocation(mock.Anything, coreentity.WorkLocationDeleteFilter{
						TenantID: "tenant-1",
						ID:       "wl-1",
					}).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkLocationDeleteFilter{
				TenantID: "tenant-1",
				ID:       "wl-1",
			},
			setup: func(repo *mocks.WorkLocationRepository, orgUnitRepo *mocks.OrgUnitRepository) {
				repo.EXPECT().
					DeleteWorkLocation(mock.Anything, coreentity.WorkLocationDeleteFilter{
						TenantID: "tenant-1",
						ID:       "wl-1",
					}).
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

			err := core.DeleteWorkLocation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
