package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckOut_CreatesAttendanceLogWhenCheckInExists(t *testing.T) {
	ctx := context.Background()
	loggedAt := "2026-04-30T09:00:00Z"
	data := coreentity.AttendanceLogAction{
		UserCtx:      attendanceUserContext("employee"),
		TenantID:     "tenant-1",
		LoggedAt:     loggedAt,
		SelfieFileID: "file-1",
	}
	want := &coreentity.AttendanceLog{ID: "log-2", TenantID: "tenant-1", Type: common.AttendanceTypeCheckOut}
	repo := dbmocks.NewAttendanceRepository(t)
	storageRepo := dbmocks.NewStorageRepository(t)
	tx := dbmocks.NewTransactor(t)
	runTransaction(t, tx)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(attendanceEmployee(), nil)
	repo.EXPECT().
		GetWorkShift(mock.Anything, coreentity.WorkShift{TenantID: "tenant-1", ID: "shift-1"}).
		Return(attendanceShift(), nil)
	repo.EXPECT().
		ExistsAttendanceByTypeOnDate(mock.Anything, "tenant-1", "employee-1", "2026-04-30", "check_out").
		Return(false, nil)
	repo.EXPECT().
		ExistsAttendanceByTypeOnDate(mock.Anything, "tenant-1", "employee-1", "2026-04-30", "check_in").
		Return(true, nil)
	storageRepo.EXPECT().
		GetFile(mock.Anything, coreentity.FileFilter{TenantID: "tenant-1", ID: "file-1"}).
		Return(&coreentity.File{
			ID:       "file-1",
			TenantID: "tenant-1",
			Folder:   common.S3FolderAttendanceFace,
			Status:   common.FileStatusPending,
		}, nil)
	repo.EXPECT().
		CreateAttendanceLog(mock.Anything, mock.MatchedBy(func(log coreentity.AttendanceLog) bool {
			return log.TenantID == "tenant-1" &&
				log.EmployeeID == "employee-1" &&
				log.AttendanceDate == "2026-04-30" &&
				log.Type == common.AttendanceTypeCheckOut &&
				log.LoggedAt == loggedAt
		})).
		Return(want, nil)
	storageRepo.EXPECT().
		CreateFileLink(mock.Anything, coreentity.CreateStorageFileLinkReq{
			TenantID:      "tenant-1",
			StorageFileID: "file-1",
			ResourceType:  "attendance_log",
			ResourceID:    "log-2",
			FieldName:     "selfie",
			SortOrder:     1,
		}).
		Return(&coreentity.StorageFileLink{ID: "link-1"}, nil)
	storageRepo.EXPECT().MarkFileAttached(mock.Anything, "tenant-1", "file-1").Return(nil)
	core := NewAttendanceCore(Config{Repo: repo, StorageRepo: storageRepo, Tx: tx})

	got, err := core.CheckOut(ctx, data)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCheckOut_RejectsWhenCheckInDoesNotExist(t *testing.T) {
	ctx := context.Background()
	data := coreentity.AttendanceLogAction{
		UserCtx:      attendanceUserContext("employee"),
		TenantID:     "tenant-1",
		LoggedAt:     "2026-04-30T09:00:00Z",
		SelfieFileID: "file-1",
	}
	repo := dbmocks.NewAttendanceRepository(t)
	tx := dbmocks.NewTransactor(t)
	runTransaction(t, tx)
	repo.EXPECT().GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").Return(attendanceEmployee(), nil)
	repo.EXPECT().
		GetWorkShift(mock.Anything, coreentity.WorkShift{TenantID: "tenant-1", ID: "shift-1"}).
		Return(attendanceShift(), nil)
	repo.EXPECT().
		ExistsAttendanceByTypeOnDate(mock.Anything, "tenant-1", "employee-1", "2026-04-30", "check_out").
		Return(false, nil)
	repo.EXPECT().
		ExistsAttendanceByTypeOnDate(mock.Anything, "tenant-1", "employee-1", "2026-04-30", "check_in").
		Return(false, nil)
	core := NewAttendanceCore(Config{Repo: repo, Tx: tx})

	got, err := core.CheckOut(ctx, data)

	require.Error(t, err)
	require.Nil(t, got)
}
