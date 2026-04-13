package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *workShiftRepo) UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:workshift:update_work_shift:UpdateWorkShift")
	defer span.End()

	query := `
		UPDATE work_shifts
		SET name = ?, timezone = ?, start_time = ?, end_time = ?, grace_period_minutes = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Name, data.Timezone, data.StartTime, data.EndTime, data.GracePeriodMinutes, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update work shift")
		return err
	}
	return nil
}
