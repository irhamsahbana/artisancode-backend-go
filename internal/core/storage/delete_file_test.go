package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorageCore_DeleteFile(t *testing.T) {
	ctx := context.Background()
	input := &coreentity.DeleteFileReq{
		TenantID: "tenant-1",
		Filename: "private/tenants/tenant-1/profile.jpg",
	}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		wantError bool
	}{
		{
			name: "success",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, coreentity.FileFilter{
						TenantID: input.TenantID,
						Filename: input.Filename,
					}).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: input.Filename,
						Status:   common.FileStatusPending,
					}, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: input.Filename}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-1").
					Return(nil)
			},
		},
		{
			name: "returns repository error when file lookup fails",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, coreentity.FileFilter{
						TenantID: input.TenantID,
						Filename: input.Filename,
					}).
					Return(nil, errors.New("file lookup failed"))
			},
			wantError: true,
		},
		{
			name: "returns business error when file is not pending",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, coreentity.FileFilter{
						TenantID: input.TenantID,
						Filename: input.Filename,
					}).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: input.Filename,
						Status:   common.FileStatusAttached,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "returns storage provider error",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, mock.AnythingOfType("coreentity.FileFilter")).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: input.Filename,
						Status:   common.FileStatusPending,
					}, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: input.Filename}).
					Return(errors.New("delete object failed"))
			},
			wantError: true,
		},
		{
			name: "returns repository error when marking deleted fails",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, mock.AnythingOfType("coreentity.FileFilter")).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: input.Filename,
						Status:   common.FileStatusPending,
					}, nil)
				s3.EXPECT().
					DeleteFile(mock.Anything, &coreentity.DeleteFileReq{Filename: input.Filename}).
					Return(nil)
				repo.EXPECT().
					MarkFileDeleted(mock.Anything, "tenant-1", "file-1").
					Return(errors.New("mark deleted failed"))
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
			err := core.DeleteFile(ctx, input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
