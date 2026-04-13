package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *storageCore) DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:storage:delete_file:DeleteFile")
	defer span.End()

	file, err := c.repo.GetFile(ctx, coreentity.FileFilter{
		TenantID: req.TenantID,
		Filename: req.Filename,
	})
	if err != nil {
		return err
	}
	if file.Status != common.FileStatusPending {
		return errmsg.NewCustomErrors(400).SetMessage("File can no longer be deleted")
	}

	err = c.s3.DeleteFile(ctx, &coreentity.DeleteFileReq{
		Filename: file.Filename,
	})
	if err != nil {
		return err
	}

	return c.repo.MarkFileDeleted(ctx, file.TenantID, file.ID)
}
