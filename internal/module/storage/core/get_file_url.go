package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *storageCore) GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetFileURL")
	defer span.End()

	file, err := c.repo.GetFile(ctx, filter)
	if err != nil {
		return "", err
	}

	return c.s3.GetFileURL(ctx, coreentity.FileFilter{
		TenantID: file.TenantID,
		Filename: file.Filename,
	})
}
