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

func TestWorkShiftCore_CreateWorkShift(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.WorkShift
		setup     func(repo *mocks.WorkShiftRepository)
		want      *coreentity.WorkShift
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.WorkShift{
				TenantID:           "tenant-1",
				Name:               "Morning Shift",
				StartTime:          "08:00",
				EndTime:            "17:00",
				Timezone:           "Asia/Jakarta",
				GracePeriodMinutes: 15,
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					CreateWorkShift(mock.Anything, mock.AnythingOfType("coreentity.WorkShift")).
					Return(&coreentity.WorkShift{
						ID:                 "shift-1",
						TenantID:           "tenant-1",
						Name:               "Morning Shift",
						StartTime:          "08:00",
						EndTime:            "17:00",
						Timezone:           "Asia/Jakarta",
						GracePeriodMinutes: 15,
					}, nil)
			},
			want: &coreentity.WorkShift{
				ID:                 "shift-1",
				TenantID:           "tenant-1",
				Name:               "Morning Shift",
				StartTime:          "08:00",
				EndTime:            "17:00",
				Timezone:           "Asia/Jakarta",
				GracePeriodMinutes: 15,
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkShift{
				TenantID:  "tenant-1",
				Name:      "Morning Shift",
				StartTime: "08:00",
				EndTime:   "17:00",
				Timezone:  "Asia/Jakarta",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					CreateWorkShift(mock.Anything, mock.AnythingOfType("coreentity.WorkShift")).
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, err := core.CreateWorkShift(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.ID, got.ID)
			require.Equal(t, tt.want.TenantID, got.TenantID)
		})
	}
}
