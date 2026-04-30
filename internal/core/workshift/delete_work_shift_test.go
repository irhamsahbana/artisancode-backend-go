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

func TestWorkShiftCore_DeleteWorkShift(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.WorkShiftDeleteFilter
		setup     func(repo *mocks.WorkShiftRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkShiftDeleteFilter{
				TenantID: "tenant-1",
				ID:       "shift-1",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					DeleteWorkShift(mock.Anything, coreentity.WorkShiftDeleteFilter{
						TenantID: "tenant-1",
						ID:       "shift-1",
					}).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkShiftDeleteFilter{
				TenantID: "tenant-1",
				ID:       "shift-1",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					DeleteWorkShift(mock.Anything, coreentity.WorkShiftDeleteFilter{
						TenantID: "tenant-1",
						ID:       "shift-1",
					}).
					Return(errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewWorkShiftRepository(t)

			if tt.setup != nil {
				tt.setup(repo)
			}

			core := NewWorkShiftCore(Config{
				Repo: repo,
			})

			err := core.DeleteWorkShift(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
