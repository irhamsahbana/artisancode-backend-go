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

func TestGetAttendancePolicy_ReturnsShiftAndLocationPolicy(t *testing.T) {
	ctx := context.Background()
	filter := attendanceSelfFilter("employee")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(attendanceEmployee(), nil)
	repo.EXPECT().
		GetWorkShift(mock.Anything, coreentity.WorkShift{TenantID: "tenant-1", ID: "shift-1"}).
		Return(attendanceShift(), nil)
	repo.EXPECT().
		GetWorkLocation(mock.Anything, coreentity.WorkLocation{TenantID: "tenant-1", ID: "location-1"}).
		Return(&coreentity.WorkLocation{ID: "location-1", RadiusMeters: intPtr(75)}, nil)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendancePolicy(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, &coreentity.AttendancePolicy{
		UserCtx:                 filter.UserCtx,
		Timezone:                "UTC",
		AttendanceRadiusMeters:  75,
		AttendanceCheckInStart:  "08:00",
		AttendanceCheckInEnd:    "08:00",
		AttendanceCheckOutStart: "17:00",
		AttendanceCheckOutEnd:   "17:00",
	}, got)
}

func TestGetAttendancePolicy_RejectsEmployeeWithoutShift(t *testing.T) {
	ctx := context.Background()
	filter := attendanceSelfFilter("employee")
	employee := attendanceEmployee()
	employee.ShiftID = nil
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(employee, nil)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendancePolicy(ctx, filter)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestGetAttendancePolicy_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := attendanceSelfFilter("employee")
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(nil, wantErr)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendancePolicy(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
