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
	ctx, span := tracing.StartSpan(ctx, "repo.MarkFileAttached")
	defer span.End()

	query := `
		UPDATE storage_files
		SET
			status = ?,
			updated_at = ?
		WHERE tenant_id = ? AND id = ? AND status = ? AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	result, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(query),
		common.FileStatusAttached,
		now,
		tenantID,
		id,
		common.FileStatusPending,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"id":        id,
		}).Msg("Failed to mark storage file attached")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"id":        id,
		}).Msg("Failed to read storage file attach result")
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("File not found")
	}

	return nil
}
