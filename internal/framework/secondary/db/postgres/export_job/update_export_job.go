package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *exportJobRepo) UpdateExportJob(ctx context.Context, data coreentity.ExportJobUpdate) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:export_job:update_export_job:UpdateExportJob",
	)
	defer span.End()

	query := `
		UPDATE export_jobs
		SET
			status = ?,
			file_id = ?,
			error_message = ?,
			started_at = COALESCE(?, started_at),
			completed_at = ?,
			expires_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = ? AND id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(
		ctx,
		exec.Rebind(query),
		data.Status,
		data.FileID,
		data.ErrorMessage,
		data.StartedAt,
		data.CompletedAt,
		data.ExpiresAt,
		data.TenantID,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update export job")
		return err
	}

	return nil
}
