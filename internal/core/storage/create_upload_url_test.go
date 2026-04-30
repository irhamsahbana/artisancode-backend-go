package core

import (
	"context"
	"errors"
	"strings"
	"testing"

	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorageCore_CreateUploadURL(t *testing.T) {
	ctx := context.Background()

	input := coreentity.PresignUploadURLReq{
		TenantID:         "tenant-1",
		CreatedBy:        "user-1",
		Filename:         "profile photo.jpg",
		OriginalFilename: strPtr("profile photo.jpg"),
		ContentType:      "image/jpeg",
		IsPublic:         false,
	}

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract)
		want      *coreentity.PresignUploadURLResp
		wantError bool
	}{
		{
			name: "success",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					CreateFile(mock.Anything, mock.MatchedBy(func(file coreentity.File) bool {
						return file.TenantID == "tenant-1" &&
							file.CreatedBy == "user-1" &&
							file.OriginalFilename != nil &&
							*file.OriginalFilename == "profile photo.jpg" &&
							strings.HasPrefix(file.Filename, "private/tenants/tenant-1/") &&
							strings.HasSuffix(file.Filename, "-profile-photo.jpg")
					})).
					Return(&coreentity.File{ID: "file-1"}, nil)
				s3.EXPECT().
					PresignUploadURL(mock.Anything, mock.MatchedBy(func(req *coreentity.PresignUploadURLReq) bool {
						return req != nil &&
							strings.HasPrefix(req.Filename, "private/tenants/tenant-1/") &&
							strings.HasSuffix(req.Filename, "-profile-photo.jpg")
					})).
					Return(&coreentity.PresignUploadURLResp{
						URL: "https://s3.example.com/upload-url",
					}, nil)
			},
			want: &coreentity.PresignUploadURLResp{
				URL:    "https://s3.example.com/upload-url",
				FileID: "file-1",
			},
		},
		{
			name: "returns repository error",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					CreateFile(mock.Anything, mock.AnythingOfType("coreentity.File")).
					Return(nil, errors.New("create file failed"))
			},
			wantError: true,
		},
		{
			name: "returns storage provider error",
			setup: func(repo *dbMocks.StorageRepository, s3 *integrationMocks.StorageContract) {
				repo.EXPECT().
					CreateFile(mock.Anything, mock.AnythingOfType("coreentity.File")).
					Return(&coreentity.File{ID: "file-1"}, nil)
				s3.EXPECT().
					PresignUploadURL(mock.Anything, mock.AnythingOfType("*coreentity.PresignUploadURLReq")).
					Return(nil, errors.New("presign failed"))
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
			got, err := core.CreateUploadURL(ctx, input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
