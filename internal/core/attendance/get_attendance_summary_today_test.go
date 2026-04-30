package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAttendanceSummaryToday_ReturnsCheckedInSummary(t *testing.T) {
	ctx := context.Background()
	filter := attendanceSelfFilter("employee")
	repo := dbmocks.NewAttendanceRepository(t)
	employee := attendanceEmployee()
	employee.LocationID = nil
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(employee, nil).Once()
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(employee, nil).Once()
	repo.EXPECT().
		GetWorkShift(mock.Anything, coreentity.WorkShift{TenantID: "tenant-1", ID: "shift-1"}).
		Return(attendanceShift(), nil)
	repo.EXPECT().
		GetAttendanceLogs(mock.Anything, mock.MatchedBy(func(got coreentity.AttendanceLogListFilter) bool {
			return got.TenantID == "tenant-1" &&
				got.EmployeeID != nil &&
				*got.EmployeeID == "employee-1" &&
				got.AttendanceDay != nil &&
				got.Page == 1 &&
				got.Paginate == 10
		})).
		Return([]coreentity.AttendanceLog{
			{
				ID:       "log-1",
				TenantID: "tenant-1",
				Type:     common.AttendanceTypeCheckIn,
				LoggedAt: "2026-04-30T01:00:00Z",
			},
		}, 1, nil)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendanceSummaryToday(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, filter.UserCtx, got.UserCtx)
	requireDateNearNow(t, got.AttendanceDate, time.UTC)
	require.True(t, got.CheckedIn)
	require.False(t, got.CheckedOut)
	require.False(t, got.CanCheckIn)
	require.True(t, got.CanCheckOut)
	require.Equal(t, common.AttendanceTodayStatusCheckedIn, got.TodayStatus)
	require.Equal(t, "log-1", *got.CheckInLogID)
	require.Nil(t, got.CheckOutLogID)
	require.Equal(t, common.AttendanceTypeCheckIn, *got.LastLogType)
	require.Equal(t, "2026-04-30T01:00:00Z", *got.LastLoggedAt)
}

func TestGetAttendanceSummaryToday_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := attendanceSelfFilter("employee")
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(nil, wantErr)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendanceSummaryToday(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
