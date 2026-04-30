package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	repositoryMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExportJobCore_GetExportJobs(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.ExportJobListFilter{
		UserCtx:  common.UserContext{UserID: "user-1", TenantID: "tenant-1", Roles: []string{"admin"}},
		TenantID: "tenant-1",
		Page:     1,
		Paginate: 10,
	}
	fileID := "file-1"
	expiresAt := "2999-01-01T00:00:00Z"
	originalFilename := "attendance-report.csv"

	tests := []struct {
		name      string
		filter    coreentity.ExportJobListFilter
		setup     func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		want      []coreentity.ExportJob
		wantTotal int
		wantError bool
	}{
		{
			name:   "success attaches download url for completed jobs",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJobs(mock.Anything, filter).
					Return([]coreentity.ExportJob{
						{
							ID:        "export-1",
							TenantID:  "tenant-1",
							Status:    coreentity.ExportJobStatusCompleted,
							FileID:    &fileID,
							ExpiresAt: &expiresAt,
						},
						{
							ID:       "export-2",
							TenantID: "tenant-1",
							Status:   coreentity.ExportJobStatusPending,
						},
					}, 2, nil)
				storageRepo.EXPECT().
					GetFile(mock.Anything, coreentity.FileFilter{
						TenantID: "tenant-1",
						ID:       "file-1",
					}).
					Return(&coreentity.File{
						ID:               "file-1",
						TenantID:         "tenant-1",
						Filename:         "private/report.csv",
						OriginalFilename: &originalFilename,
					}, nil)
				s3.EXPECT().
					GetFileURL(mock.Anything, coreentity.FileFilter{
						TenantID: "tenant-1",
						Filename: "private/report.csv",
					}).
					Return("https://signed.example.com/report.csv", nil)
			},
			want: []coreentity.ExportJob{
				{
					ID:          "export-1",
					TenantID:    "tenant-1",
					Status:      coreentity.ExportJobStatusCompleted,
					FileID:      &fileID,
					FileName:    &originalFilename,
					ExpiresAt:   &expiresAt,
					DownloadURL: stringPtr("https://signed.example.com/report.csv"),
				},
				{
					ID:       "export-2",
					TenantID: "tenant-1",
					Status:   coreentity.ExportJobStatusPending,
				},
			},
			wantTotal: 2,
		},
		{
			name: "returns authorization error",
			filter: coreentity.ExportJobListFilter{
				UserCtx:  common.UserContext{UserID: "user-2", TenantID: "tenant-1", Roles: []string{"member"}},
				TenantID: "tenant-1",
			},
			wantError: true,
		},
		{
			name:   "returns repository error",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJobs(mock.Anything, filter).
					Return(nil, 0, errors.New("list export jobs failed"))
			},
			wantError: true,
		},
		{
			name:   "returns signing error",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJobs(mock.Anything, filter).
					Return([]coreentity.ExportJob{
						{
							ID:        "export-1",
							TenantID:  "tenant-1",
							Status:    coreentity.ExportJobStatusCompleted,
							FileID:    &fileID,
							ExpiresAt: &expiresAt,
						},
					}, 1, nil)
				storageRepo.EXPECT().
					GetFile(mock.Anything, mock.AnythingOfType("coreentity.FileFilter")).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: "private/report.csv",
					}, nil)
				s3.EXPECT().
					GetFileURL(mock.Anything, mock.AnythingOfType("coreentity.FileFilter")).
					Return("", errors.New("sign url failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositoryMocks.NewExportJobRepository(t)
			storageRepo := repositoryMocks.NewStorageRepository(t)
			s3 := integrationMocks.NewStorageContract(t)
			if tt.setup != nil {
				tt.setup(repo, storageRepo, s3)
			}

			core := NewExportJobCore(Config{
				Repo:        repo,
				StorageRepo: storageRepo,
				S3:          s3,
			})
			got, total, err := core.GetExportJobs(ctx, tt.filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantTotal, total)
			require.Equal(t, tt.want, got)
		})
	}
}
