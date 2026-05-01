package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) MarkFileAttached(ctx context.Context, tenantID, id string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:storage:mark_file_attached:MarkFileAttached",
	)
	defer span.End()

	query := `
		UPDATE storage_files
		SET
			status = ?,
			updated_at = ?
		WHERE tenant_id = ? AND id = ? AND status = ? AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	exec := r.executor(ctx)
	payload := map[string]string{
		"tenant_id": tenantID,
		"id":        id,
	}
	result, err := exec.ExecContext(
		ctx,
		exec.Rebind(query),
		common.FileStatusAttached,
		now,
		tenantID,
		id,
		common.FileStatusPending,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to mark storage file attached")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to read storage file attach result")
		return err
	}
	if rowsAffected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("File not found when marking attached")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageFileNotFound)
	}

	return nil
}
