package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *storageCore) UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:storage:upload_file:UploadFile")
	defer span.End()

	return c.s3.UploadFile(ctx, req)
}
