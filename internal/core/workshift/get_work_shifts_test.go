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

func TestWorkShiftCore_GetWorkShifts(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    coreentity.WorkShiftListFilter
		setup     func(repo *mocks.WorkShiftRepository)
		want      []coreentity.WorkShift
		wantCount int
		wantError bool
	}{
		{
			name: "success",
			filter: coreentity.WorkShiftListFilter{
				TenantID: "tenant-1",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					GetWorkShifts(mock.Anything, coreentity.WorkShiftListFilter{
						TenantID: "tenant-1",
						Page:     1,
						Paginate: 10,
					}).
					Return([]coreentity.WorkShift{
						{
							ID:        "shift-1",
							TenantID:  "tenant-1",
							Name:      "Morning Shift",
							StartTime: "08:00",
							EndTime:   "17:00",
						},
					}, 1, nil)
			},
			want: []coreentity.WorkShift{
				{
					ID:        "shift-1",
					TenantID:  "tenant-1",
					Name:      "Morning Shift",
					StartTime: "08:00",
					EndTime:   "17:00",
				},
			},
			wantCount: 1,
		},
		{
			name: "dependency error from repository",
			filter: coreentity.WorkShiftListFilter{
				TenantID: "tenant-1",
				Page:     1,
				Paginate: 10,
			},
			setup: func(repo *mocks.WorkShiftRepository) {
				repo.EXPECT().
					GetWorkShifts(mock.Anything, mock.AnythingOfType("coreentity.WorkShiftListFilter")).
					Return(nil, 0, errmsg.NewCustomErrors(500).SetMessage("database error"))
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

			got, count, err := core.GetWorkShifts(ctx, tt.filter)

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
