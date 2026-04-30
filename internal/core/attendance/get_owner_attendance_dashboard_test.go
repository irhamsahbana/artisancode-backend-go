package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOwnerAttendanceDashboard_ReturnsDashboardForOwner(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OwnerAttendanceDashboardFilter{
		UserCtx:  attendanceUserContext("owner"),
		TenantID: "tenant-1",
		Date:     "2026-04-30",
	}
	summary := &coreentity.OwnerAttendanceDashboardSummary{ActiveEmployeeCount: 3}
	exceptions := []coreentity.OwnerAttendanceDashboardException{{EmployeeID: "employee-1"}}
	trend := []coreentity.OwnerAttendanceDashboardDailyTrend{{AttendanceDate: "2026-04-30"}}
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetOwnerAttendanceDashboardSummary(mock.Anything, filter).Return(summary, nil)
	repo.EXPECT().GetOwnerAttendanceDashboardExceptions(mock.Anything, filter).Return(exceptions, nil)
	repo.EXPECT().GetOwnerAttendanceDashboardTrend(mock.Anything, filter).Return(trend, nil)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetOwnerAttendanceDashboard(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, filter.UserCtx, got.UserCtx)
	require.Equal(t, "2026-04-30", got.AttendanceDate)
	require.Equal(t, *summary, got.Summary)
	require.Equal(t, exceptions, got.TodayExceptions)
	require.Equal(t, trend, got.DailyTrend)
}

func TestGetOwnerAttendanceDashboard_AllowsAdminRole(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OwnerAttendanceDashboardFilter{
		UserCtx:  attendanceUserContext("admin"),
		TenantID: "tenant-1",
		Date:     "2026-04-30",
	}
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().
		GetOwnerAttendanceDashboardSummary(mock.Anything, filter).
		Return(&coreentity.OwnerAttendanceDashboardSummary{}, nil)
	repo.EXPECT().
		GetOwnerAttendanceDashboardExceptions(mock.Anything, filter).
		Return([]coreentity.OwnerAttendanceDashboardException{}, nil)
	repo.EXPECT().
		GetOwnerAttendanceDashboardTrend(mock.Anything, filter).
		Return([]coreentity.OwnerAttendanceDashboardDailyTrend{}, nil)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetOwnerAttendanceDashboard(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, got)
}

func TestGetOwnerAttendanceDashboard_RejectsNonOwnerAndNonAdmin(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OwnerAttendanceDashboardFilter{
		UserCtx:  attendanceUserContext("employee"),
		TenantID: "tenant-1",
		Date:     "2026-04-30",
	}
	core := NewAttendanceCore(Config{})

	got, err := core.GetOwnerAttendanceDashboard(ctx, filter)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestGetOwnerAttendanceDashboard_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OwnerAttendanceDashboardFilter{
		UserCtx:  attendanceUserContext("owner"),
		TenantID: "tenant-1",
		Date:     "2026-04-30",
	}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetOwnerAttendanceDashboardSummary(mock.Anything, filter).Return(nil, wantErr)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetOwnerAttendanceDashboard(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
