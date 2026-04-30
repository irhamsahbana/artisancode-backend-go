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

func TestWorkShiftCore_UpdateWorkShift(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.WorkShift
		setup     func(repo *mocks.WorkShiftRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkShift{
				TenantID:  "tenant-1",
				ID:        "shift-1",
				Name:      "Updated Shift",
				StartTime: "09:00",
				EndTime:   "18:00",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					UpdateWorkShift(mock.Anything, mock.AnythingOfType("coreentity.WorkShift")).
					Return(nil)
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkShift{
				TenantID:  "tenant-1",
				ID:        "shift-1",
				Name:      "Updated Shift",
				StartTime: "09:00",
				EndTime:   "18:00",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					UpdateWorkShift(mock.Anything, mock.AnythingOfType("coreentity.WorkShift")).
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

			err := core.UpdateWorkShift(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
