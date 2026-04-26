package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) MarkFileDeleted(ctx context.Context, tenantID, id string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:storage:mark_file_deleted:MarkFileDeleted",
	)
	defer span.End()

	query := `
		UPDATE storage_files
		SET
			status = ?,
			deleted_at = ?,
			updated_at = ?
		WHERE tenant_id = ? AND id = ? AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	exec := r.executor(ctx)
	_, err := exec.ExecContext(
		ctx,
		exec.Rebind(query),
		common.FileStatusDeleted,
		now,
		now,
		tenantID,
		id,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"id":        id,
		}).Msg("Failed to mark storage file deleted")
		return err
	}

	return nil
}
