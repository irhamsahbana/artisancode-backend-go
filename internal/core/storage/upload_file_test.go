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

func TestStorageCore_UploadFile(t *testing.T) {
	ctx := context.Background()
	req := &coreentity.UploadFileReq{
		Filename: "profile.jpg",
		TenantID: "tenant-1",
	}

	tests := []struct {
		name      string
		setup     func(s3 *integrationMocks.StorageContract)
		want      *coreentity.UploadFileResp
		wantError bool
	}{
		{
			name: "success",
			setup: func(s3 *integrationMocks.StorageContract) {
				s3.EXPECT().
					UploadFile(mock.Anything, req).
					Return(&coreentity.UploadFileResp{
						Filename: "profile.jpg",
						URL:      "https://cdn.example.com/profile.jpg",
					}, nil)
			},
			want: &coreentity.UploadFileResp{
				Filename: "profile.jpg",
				URL:      "https://cdn.example.com/profile.jpg",
			},
		},
		{
			name: "returns storage provider error",
			setup: func(s3 *integrationMocks.StorageContract) {
				s3.EXPECT().
					UploadFile(mock.Anything, req).
					Return(nil, errors.New("upload failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewStorageRepository(t)
			s3 := integrationMocks.NewStorageContract(t)
			tt.setup(s3)

			core := NewStorageCore(s3, repo)
			got, err := core.UploadFile(ctx, req)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
