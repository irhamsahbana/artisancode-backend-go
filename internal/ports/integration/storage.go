package integration

import (
	"codebase-app/internal/entity/coreentity"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type StorageContract interface {
	UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error)
	DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error
	ListFiles(ctx context.Context) ([]types.Object, error)
	GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error)
}
