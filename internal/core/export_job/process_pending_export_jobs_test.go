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

func TestExportJobCore_ProcessPendingExportJobs(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		limit int
		setup func(
			repo *repositoryMocks.ExportJobRepository,
			attendanceRepo *repositoryMocks.AttendanceRepository,
			storageRepo *repositoryMocks.StorageRepository,
			s3 *integrationMocks.StorageContract,
		)
		want      *coreentity.ExportJobProcessResult
		wantError bool
	}{
		{
			name:  "returns empty result when no job is claimed",
			limit: 10,
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					ClaimPendingExportJob(mock.Anything).
					Return(nil, nil)
			},
			want: &coreentity.ExportJobProcessResult{},
		},
		{
			name:  "uses minimum limit and processes one claimed job",
			limit: 0,
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				startedAt := "2026-04-30T10:00:00Z"
				repo.EXPECT().
					ClaimPendingExportJob(mock.Anything).
					Return(&coreentity.ExportJob{
						ID:           "export-1",
						TenantID:     "tenant-1",
						RequestedBy:  "user-1",
						ProcessorKey: "attendance_logs",
						Format:       coreentity.ExportJobFormatXLSX,
						Status:       coreentity.ExportJobStatusProcessing,
						ParamsJSON:   `{"language":"en"}`,
						StartedAt:    &startedAt,
					}, nil)
				attendanceRepo.EXPECT().
					GetAttendanceLogsAll(mock.Anything, mock.MatchedBy(func(filter coreentity.AttendanceLogListFilter) bool {
						return filter.TenantID == "tenant-1" &&
							filter.Page == 1 &&
							filter.Paginate == 1000000
					})).
					Return(nil, nil)
				storageRepo.EXPECT().
					CreateFile(mock.Anything, mock.AnythingOfType("coreentity.File")).
					Return(&coreentity.File{ID: "file-1"}, nil)
				s3.EXPECT().
					UploadBytes(mock.Anything, mock.MatchedBy(func(req *coreentity.UploadBytesReq) bool {
						return req != nil &&
							req.ContentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
							len(req.Body) > 0
					})).
					Return(&coreentity.UploadFileResp{Filename: "private/report.xlsx"}, nil)
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
						return update.Status == coreentity.ExportJobStatusCompleted &&
							update.FileID != nil &&
							*update.FileID == "file-1" &&
							update.StartedAt != nil &&
							*update.StartedAt == startedAt
					})).
					Return(nil)
			},
			want: &coreentity.ExportJobProcessResult{Processed: 1},
		},
		{
			name:  "returns claim error",
			limit: 1,
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				repo.EXPECT().
					ClaimPendingExportJob(mock.Anything).
					Return(nil, errors.New("claim failed"))
			},
			wantError: true,
		},
		{
			name:  "marks job failed when processing fails",
			limit: 1,
			setup: func(
				repo *repositoryMocks.ExportJobRepository,
				attendanceRepo *repositoryMocks.AttendanceRepository,
				storageRepo *repositoryMocks.StorageRepository,
				s3 *integrationMocks.StorageContract,
			) {
				startedAt := "2026-04-30T10:00:00Z"
				repo.EXPECT().
					ClaimPendingExportJob(mock.Anything).
					Return(&coreentity.ExportJob{
						ID:           "export-1",
						TenantID:     "tenant-1",
						ProcessorKey: "unsupported",
						Status:       coreentity.ExportJobStatusProcessing,
						StartedAt:    &startedAt,
					}, nil)
				repo.EXPECT().
					UpdateExportJob(mock.Anything, mock.MatchedBy(func(update coreentity.ExportJobUpdate) bool {
						return update.TenantID == "tenant-1" &&
							update.ID == "export-1" &&
							update.Status == coreentity.ExportJobStatusFailed &&
							update.ErrorMessage != nil &&
							update.StartedAt != nil &&
							*update.StartedAt == startedAt
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
			got, err := core.ProcessPendingExportJobs(ctx, tt.limit)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
