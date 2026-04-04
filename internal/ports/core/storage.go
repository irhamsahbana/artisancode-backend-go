package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type StorageCore interface {
	UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error)
	CreateUploadURL(ctx context.Context, req coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error)
	DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error
	ListFiles(ctx context.Context) ([]types.Object, error)
	GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error)
	CleanupExpiredFiles(ctx context.Context, req coreentity.CleanupExpiredFilesReq) (*coreentity.CleanupExpiredFilesResp, error)
}
