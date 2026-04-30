package core

import (
	"context"
	"errors"
	"testing"

	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorageCore_ListFiles(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(s3 *integrationMocks.StorageContract)
		want      []types.Object
		wantError bool
	}{
		{
			name: "success",
			setup: func(s3 *integrationMocks.StorageContract) {
				s3.EXPECT().
					ListFiles(mock.Anything).
					Return([]types.Object{{Key: aws.String("private/file-1.jpg")}}, nil)
			},
			want: []types.Object{{Key: aws.String("private/file-1.jpg")}},
		},
		{
			name: "returns storage provider error",
			setup: func(s3 *integrationMocks.StorageContract) {
				s3.EXPECT().
					ListFiles(mock.Anything).
					Return(nil, errors.New("list failed"))
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
			got, err := core.ListFiles(ctx)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
