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

func TestStorageCore_GetFileURL(t *testing.T) {
	ctx := context.Background()
	input := coreentity.FileFilter{
		TenantID: "tenant-1",
		ID:       "file-1",
	}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		want      string
		wantError bool
	}{
		{
			name: "success",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, input).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: "private/tenants/tenant-1/profile.jpg",
					}, nil)
				s3.EXPECT().
					GetFileURL(mock.Anything, coreentity.FileFilter{
						TenantID: "tenant-1",
						Filename: "private/tenants/tenant-1/profile.jpg",
					}).
					Return("https://signed.example.com/profile.jpg", nil)
			},
			want: "https://signed.example.com/profile.jpg",
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, input).
					Return(nil, errors.New("file not found"))
			},
			wantError: true,
		},
		{
			name: "returns storage provider error",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					GetFile(mock.Anything, input).
					Return(&coreentity.File{
						ID:       "file-1",
						TenantID: "tenant-1",
						Filename: "private/tenants/tenant-1/profile.jpg",
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
			repo := dbMocks.NewStorageRepository(t)
			s3 := integrationMocks.NewStorageContract(t)
			tt.setup(repo, s3)

			core := NewStorageCore(s3, repo)
			got, err := core.GetFileURL(ctx, input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
