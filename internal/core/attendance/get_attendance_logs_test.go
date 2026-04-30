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

func TestGetAttendanceLogs_ScopesNonOwnerToOwnEmployee(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.AttendanceLogListFilter{
		UserCtx:  attendanceUserContext("employee"),
		TenantID: "tenant-1",
		Page:     1,
		Paginate: 10,
	}
	items := []coreentity.AttendanceLog{{ID: "log-1", TenantID: "tenant-1", EmployeeID: "employee-1"}}
	repo := dbmocks.NewAttendanceRepository(t)
	storageRepo := dbmocks.NewStorageRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(attendanceEmployee(), nil)
	repo.EXPECT().
		GetAttendanceLogs(mock.Anything, mock.MatchedBy(func(got coreentity.AttendanceLogListFilter) bool {
			return got.TenantID == "tenant-1" &&
				got.EmployeeID != nil &&
				*got.EmployeeID == "employee-1"
		})).
		Return(items, 1, nil)
	storageRepo.EXPECT().
		GetFileLink(mock.Anything, coreentity.StorageFileLinkFilter{
			TenantID:     "tenant-1",
			ResourceType: "attendance_log",
			ResourceID:   "log-1",
			FieldName:    "selfie",
		}).
		Return(nil, nil)
	core := NewAttendanceCore(Config{Repo: repo, StorageRepo: storageRepo})

	got, total, err := core.GetAttendanceLogs(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, items, got)
	require.Equal(t, 1, total)
}

func TestGetAttendanceLogs_ReturnsEmptyWhenNonOwnerEmployeeLookupFails(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.AttendanceLogListFilter{
		UserCtx:  attendanceUserContext("employee"),
		TenantID: "tenant-1",
	}
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().
		GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").
		Return(nil, errors.New("repository failed"))
	core := NewAttendanceCore(Config{Repo: repo})

	got, total, err := core.GetAttendanceLogs(ctx, filter)

	require.NoError(t, err)
	require.Empty(t, got)
	require.Zero(t, total)
}

func TestGetAttendanceLogs_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.AttendanceLogListFilter{
		UserCtx:  attendanceUserContext("owner"),
		TenantID: "tenant-1",
	}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetAttendanceLogs(mock.Anything, filter).Return(nil, 0, wantErr)
	core := NewAttendanceCore(Config{Repo: repo})

	got, total, err := core.GetAttendanceLogs(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
	require.Zero(t, total)
}
