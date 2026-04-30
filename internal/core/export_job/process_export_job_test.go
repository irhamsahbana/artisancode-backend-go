package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	repositoryMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExportJobCore_ProcessExportJob(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.ExportJobDetailFilter{
		TenantID: "tenant-1",
		ID:       "export-1",
	}

	tests := []struct {
		name  string
		setup func(
			repo *repositoryMocks.ExportJobRepository,
			attendanceRepo *repositoryMocks.AttendanceRepository,
			storageRepo *repositoryMocks.StorageRepository,
			s3 *integrationMocks.StorageContract,
		)
		wantError bool
	}{
		{
			name: "success processes pending attendance export",
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:           "export-1",
						TenantID:     "tenant-1",
						RequestedBy:  "user-1",
						ProcessorKey: "attendance_logs",
						Format:       coreentity.ExportJobFormatCSV,
						Status:       coreentity.ExportJobStatusPending,
						ParamsJSON:   `{"language":"en"}`,
					}, nil)
				repo.EXPECT().
					UpdateExportJob(mock.Anything, mock.MatchedBy(func(update coreentity.ExportJobUpdate) bool {
						return update.TenantID == "tenant-1" &&
							update.ID == "export-1" &&
							update.Status == coreentity.ExportJobStatusProcessing &&
							update.StartedAt != nil
					})).
					Return(nil)
				attendanceRepo.EXPECT().
					GetAttendanceLogsAll(mock.Anything, mock.MatchedBy(func(filter coreentity.AttendanceLogListFilter) bool {
						return filter.TenantID == "tenant-1" &&
							filter.Page == 1 &&
							filter.Paginate == 1000000
					})).
					Return(nil, nil)
				storageRepo.EXPECT().
					CreateFile(mock.Anything, mock.MatchedBy(func(file coreentity.File) bool {
						return file.TenantID == "tenant-1" &&
							file.CreatedBy == "user-1" &&
							file.OriginalFilename != nil &&
							file.ContentType != nil &&
							*file.ContentType == "text/csv"
					})).
					Return(&coreentity.File{ID: "file-1"}, nil)
				s3.EXPECT().
					UploadBytes(mock.Anything, mock.MatchedBy(func(req *coreentity.UploadBytesReq) bool {
						return req != nil &&
							req.ContentType == "text/csv" &&
							len(req.Body) > 0
					})).
					Return(&coreentity.UploadFileResp{Filename: "private/report.csv"}, nil)
				storageRepo.EXPECT().
					CreateFileLink(mock.Anything, coreentity.CreateStorageFileLinkReq{
						TenantID:      "tenant-1",
						StorageFileID: "file-1",
						ResourceType:  "export_job",
						ResourceID:    "export-1",
						FieldName:     "report_file",
						SortOrder:     1,
					}).
					Return(&coreentity.StorageFileLink{ID: "link-1"}, nil)
				storageRepo.EXPECT().
					MarkFileAttached(mock.Anything, "tenant-1", "file-1").
					Return(nil)
				repo.EXPECT().
					UpdateExportJob(mock.Anything, mock.MatchedBy(func(update coreentity.ExportJobUpdate) bool {
						return update.TenantID == "tenant-1" &&
							update.ID == "export-1" &&
							update.Status == coreentity.ExportJobStatusCompleted &&
							update.FileID != nil &&
							*update.FileID == "file-1" &&
							update.StartedAt != nil &&
							update.CompletedAt != nil &&
							update.ExpiresAt != nil
					})).
					Return(nil)
			},
		},
		{
			name: "does nothing when job is not pending",
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:       "export-1",
						TenantID: "tenant-1",
						Status:   coreentity.ExportJobStatusCompleted,
					}, nil)
			},
		},
		{
			name: "returns repository lookup error",
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(nil, errors.New("get export job failed"))
			},
			wantError: true,
		},
		{
			name: "returns processing status update error",
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:           "export-1",
						TenantID:     "tenant-1",
						ProcessorKey: "attendance_logs",
						Status:       coreentity.ExportJobStatusPending,
						ParamsJSON:   `{}`,
					}, nil)
				repo.EXPECT().
					UpdateExportJob(mock.Anything, mock.MatchedBy(func(update coreentity.ExportJobUpdate) bool {
						return update.Status == coreentity.ExportJobStatusProcessing
					})).
					Return(errors.New("update processing failed"))
			},
			wantError: true,
		},
		{
			name: "returns processor error",
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:           "export-1",
						TenantID:     "tenant-1",
						ProcessorKey: "unknown",
						Status:       coreentity.ExportJobStatusPending,
					}, nil)
				repo.EXPECT().
					UpdateExportJob(mock.Anything, mock.MatchedBy(func(update coreentity.ExportJobUpdate) bool {
						return update.Status == coreentity.ExportJobStatusProcessing
					})).
					Return(nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositoryMocks.NewExportJobRepository(t)
			attendanceRepo := repositoryMocks.NewAttendanceRepository(t)
			storageRepo := repositoryMocks.NewStorageRepository(t)
			s3 := integrationMocks.NewStorageContract(t)
			tt.setup(repo, attendanceRepo, storageRepo, s3)

			core := NewExportJobCore(Config{
				Repo:           repo,
				AttendanceRepo: attendanceRepo,
				StorageRepo:    storageRepo,
				S3:             s3,
			})
			err := core.ProcessExportJob(ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
