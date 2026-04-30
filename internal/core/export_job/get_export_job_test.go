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

func TestExportJobCore_GetExportJob(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.ExportJobDetailFilter{
		UserCtx:  common.UserContext{UserID: "user-1", TenantID: "tenant-1", Roles: []string{"owner"}},
		TenantID: "tenant-1",
		ID:       "export-1",
	}
	fileID := "file-1"
	expiresAt := "2999-01-01T00:00:00Z"
	originalFilename := "attendance-report.csv"

	tests := []struct {
		name      string
		filter    coreentity.ExportJobDetailFilter
		setup     func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		want      *coreentity.ExportJob
		wantError bool
	}{
		{
			name:   "success attaches download url for completed job",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:        "export-1",
						TenantID:  "tenant-1",
						Status:    coreentity.ExportJobStatusCompleted,
						FileID:    &fileID,
						ExpiresAt: &expiresAt,
					}, nil)
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
			want: &coreentity.ExportJob{
				ID:          "export-1",
				TenantID:    "tenant-1",
				Status:      coreentity.ExportJobStatusCompleted,
				FileID:      &fileID,
				FileName:    &originalFilename,
				ExpiresAt:   &expiresAt,
				DownloadURL: stringPtr("https://signed.example.com/report.csv"),
			},
		},
		{
			name: "returns authorization error",
			filter: coreentity.ExportJobDetailFilter{
				UserCtx:  common.UserContext{UserID: "user-2", TenantID: "tenant-1", Roles: []string{"member"}},
				TenantID: "tenant-1",
				ID:       "export-1",
			},
			wantError: true,
		},
		{
			name:   "returns repository error",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(nil, errors.New("get export job failed"))
			},
			wantError: true,
		},
		{
			name:   "marks completed job expired without signing when expiration passed",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				expiredAt := "2000-01-01T00:00:00Z"
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:        "export-1",
						TenantID:  "tenant-1",
						Status:    coreentity.ExportJobStatusCompleted,
						FileID:    &fileID,
						ExpiresAt: &expiredAt,
					}, nil)
			},
			want: &coreentity.ExportJob{
				ID:        "export-1",
				TenantID:  "tenant-1",
				Status:    coreentity.ExportJobStatusExpired,
				FileID:    &fileID,
				ExpiresAt: stringPtr("2000-01-01T00:00:00Z"),
			},
		},
		{
			name:   "returns storage repository error while attaching download url",
			filter: filter,
			setup: func(repo *repositoryMocks.ExportJobRepository, storageRepo *repositoryMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExportJob(mock.Anything, filter).
					Return(&coreentity.ExportJob{
						ID:        "export-1",
						TenantID:  "tenant-1",
						Status:    coreentity.ExportJobStatusCompleted,
						FileID:    &fileID,
						ExpiresAt: &expiresAt,
					}, nil)
				storageRepo.EXPECT().
					GetFile(mock.Anything, mock.AnythingOfType("coreentity.FileFilter")).
					Return(nil, errors.New("file lookup failed"))
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
			got, err := core.GetExportJob(ctx, tt.filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
