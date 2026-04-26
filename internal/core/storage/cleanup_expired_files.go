package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (c *storageCore) CleanupExpiredFiles(
	ctx context.Context,
	req coreentity.CleanupExpiredFilesReq,
) (*coreentity.CleanupExpiredFilesResp, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:storage:cleanup_expired_files:CleanupExpiredFiles")
	defer span.End()

	before := req.Before
	if before == "" {
		before = time.Now().UTC().Format(time.RFC3339)
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}

	files, err := c.repo.GetExpiredPendingFiles(ctx, coreentity.ExpiredPendingFileFilter{
		Before: before,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp := &coreentity.CleanupExpiredFilesResp{
		Scanned: len(files),
		Deleted: 0,
	}

	for _, file := range files {
		err = c.s3.DeleteFile(ctx, &coreentity.DeleteFileReq{
			Filename: file.Filename,
		})
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any("file_id", file.ID).
				Any("object_key", file.Filename).
				Msg("Failed to delete expired storage file from object storage")
			continue
		}

		err = c.repo.MarkFileDeleted(ctx, file.TenantID, file.ID)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any("file_id", file.ID).
				Any("object_key", file.Filename).
				Msg("Failed to mark expired storage file deleted")
			continue
		}

		resp.Deleted++
	}

	return resp, nil
}
