package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/internal/ports/integration"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type storageCore struct {
	s3 integration.StorageContract
}

var _ corePorts.StorageCore = &storageCore{}

func NewStorageCore(s3 integration.StorageContract) *storageCore {
	return &storageCore{
		s3: s3,
	}
}

func (c *storageCore) UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
	ctx, span := tracing.StartSpan(ctx, "core.UploadFile")
	defer span.End()

	return c.s3.UploadFile(ctx, req)
}

func (c *storageCore) DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteFile")
	defer span.End()

	return c.s3.DeleteFile(ctx, req)
}

func (c *storageCore) ListFiles(ctx context.Context) ([]types.Object, error) {
	ctx, span := tracing.StartSpan(ctx, "core.ListFiles")
	defer span.End()

	return c.s3.ListFiles(ctx)
}

func (c *storageCore) GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetFileURL")
	defer span.End()

	return c.s3.GetFileURL(ctx, filter)
}
