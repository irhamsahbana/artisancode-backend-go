package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorageCore_CleanupExpiredFiles(t *testing.T) {
	ctx := context.Background()
	req := coreentity.CleanupExpiredFilesReq{
		Before: "2026-04-30T10:00:00Z",
		Limit:  2,
	}

	tests := []struct {
		name      string
		req       coreentity.CleanupExpiredFilesReq
		setup     func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		want      *coreentity.CleanupExpiredFilesResp
		wantError bool
	}{
		{
			name: "success",
			req:  req,
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				files := []coreentity.File{
					{ID: "file-1", TenantID: "tenant-1", Filename: "private/file-1.jpg"},
					{ID: "file-2", TenantID: "tenant-1", Filename: "private/file-2.jpg"},
				}
				repo.EXPECT().
					GetExpiredPendingFiles(mock.Anything, coreentity.ExpiredPendingFileFilter{
						Before: req.Before,
						Limit:  req.Limit,
					}).
					Return(files, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: "private/file-1.jpg"}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-1").
					Return(nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: "private/file-2.jpg"}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-2").
					Return(nil)
			},
			want: &coreentity.CleanupExpiredFilesResp{
				Scanned: 2,
				Deleted: 2,
			},
		},
		{
			name: "uses default limit when limit is invalid",
			req: coreentity.CleanupExpiredFilesReq{
				Before: req.Before,
				Limit:  0,
			},
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExpiredPendingFiles(mock.Anything, coreentity.ExpiredPendingFileFilter{
						Before: req.Before,
						Limit:  100,
					}).
					Return(nil, nil)
			},
			want: &coreentity.CleanupExpiredFilesResp{},
		},
		{
			name: "continues when object delete fails",
			req:  req,
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				files := []coreentity.File{
					{ID: "file-1", TenantID: "tenant-1", Filename: "private/file-1.jpg"},
					{ID: "file-2", TenantID: "tenant-1", Filename: "private/file-2.jpg"},
				}
				repo.EXPECT().
					GetExpiredPendingFiles(mock.Anything, mock.AnythingOfType("coreentity.ExpiredPendingFileFilter")).
					Return(files, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: "private/file-1.jpg"}).
					Return(errors.New("delete object failed"))
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: "private/file-2.jpg"}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-2").
					Return(nil)
			},
			want: &coreentity.CleanupExpiredFilesResp{
				Scanned: 2,
				Deleted: 1,
			},
		},
		{
			name: "continues when mark deleted fails",
			req:  req,
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				files := []coreentity.File{
					{ID: "file-1", TenantID: "tenant-1", Filename: "private/file-1.jpg"},
				}
				repo.EXPECT().
					GetExpiredPendingFiles(mock.Anything, mock.AnythingOfType("coreentity.ExpiredPendingFileFilter")).
					Return(files, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: "private/file-1.jpg"}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-1").
					Return(errors.New("mark deleted failed"))
			},
			want: &coreentity.CleanupExpiredFilesResp{
				Scanned: 1,
				Deleted: 0,
			},
		},
		{
			name: "returns repository scan error",
			req:  req,
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetExpiredPendingFiles(mock.Anything, mock.AnythingOfType("coreentity.ExpiredPendingFileFilter")).
					Return(nil, errors.New("scan failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewStorageRepository(t)
			s3 := integrationMocks.NewStorageContract(t)
			tt.setup(repo, s3)

			core := NewStorageCore(s3, repo)
			got, err := core.CleanupExpiredFiles(ctx, tt.req)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
