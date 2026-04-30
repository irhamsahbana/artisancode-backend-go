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

func TestWorkShiftCore_GetWorkShift(t *testing.T) {
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
				TenantID: "tenant-1",
				ID:       "shift-1",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					GetWorkShift(mock.Anything, coreentity.WorkShift{
						TenantID: "tenant-1",
						ID:       "shift-1",
					}).
					Return(&coreentity.WorkShift{
						ID:        "shift-1",
						TenantID:  "tenant-1",
						Name:      "Morning Shift",
						StartTime: "08:00",
						EndTime:   "17:00",
						Timezone:  "Asia/Jakarta",
					}, nil)
			},
			want: &coreentity.WorkShift{
				ID:        "shift-1",
				TenantID:  "tenant-1",
				Name:      "Morning Shift",
				StartTime: "08:00",
				EndTime:   "17:00",
				Timezone:  "Asia/Jakarta",
			},
		},
		{
			name: "dependency error from repository",
			input: coreentity.WorkShift{
				TenantID: "tenant-1",
				ID:       "shift-1",
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					GetWorkShift(mock.Anything, coreentity.WorkShift{
						TenantID: "tenant-1",
						ID:       "shift-1",
					}).
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

			got, err := core.GetWorkShift(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
