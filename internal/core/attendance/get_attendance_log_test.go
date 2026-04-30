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

func TestGetAttendanceLog_ScopesNonOwnerToOwnEmployeeAndAttachesNoSelfieWhenMissing(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.AttendanceLogDetailFilter{
		UserCtx:  attendanceUserContext("employee"),
		TenantID: "tenant-1",
		ID:       "log-1",
	}
	want := &coreentity.AttendanceLog{ID: "log-1", TenantID: "tenant-1", EmployeeID: "employee-1"}
	repo := dbmocks.NewAttendanceRepository(t)
	storageRepo := dbmocks.NewStorageRepository(t)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(attendanceEmployee(), nil)
	repo.EXPECT().
		GetAttendanceLog(mock.Anything, mock.MatchedBy(func(got coreentity.AttendanceLogDetailFilter) bool {
			return got.TenantID == "tenant-1" &&
				got.ID == "log-1" &&
				got.EmployeeID != nil &&
				*got.EmployeeID == "employee-1"
		})).
		Return(want, nil)
	storageRepo.EXPECT().
		GetFileLink(mock.Anything, coreentity.StorageFileLinkFilter{
			TenantID:     "tenant-1",
			ResourceType: "attendance_log",
			ResourceID:   "log-1",
			FieldName:    "selfie",
		}).
		Return(nil, nil)
	core := NewAttendanceCore(Config{Repo: repo, StorageRepo: storageRepo})

	got, err := core.GetAttendanceLog(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, want, got)
	require.Nil(t, got.SelfieFileID)
	require.Nil(t, got.SelfieURL)
}

func TestGetAttendanceLog_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.AttendanceLogDetailFilter{
		UserCtx:  attendanceUserContext("owner"),
		TenantID: "tenant-1",
		ID:       "log-1",
	}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewAttendanceRepository(t)
	repo.EXPECT().GetAttendanceLog(mock.Anything, filter).Return(nil, wantErr)
	core := NewAttendanceCore(Config{Repo: repo})

	got, err := core.GetAttendanceLog(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
