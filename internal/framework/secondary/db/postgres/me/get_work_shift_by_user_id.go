package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *meRepo) GetWorkShiftByUserID(ctx context.Context, tenantID, userID string) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:me:get_work_shift_by_user_id:GetWorkShiftByUserID")
	defer span.End()

	var data struct {
		ID                 string `db:"id"`
		TenantID           string `db:"tenant_id"`
		Name               string `db:"name"`
		Timezone           string `db:"timezone"`
		StartTime          string `db:"start_time"`
		EndTime            string `db:"end_time"`
		GracePeriodMinutes int    `db:"grace_period_minutes"`
	}

	query := `
		SELECT
			ws.id,
			ws.tenant_id,
			ws.name,
			ws.timezone,
			ws.start_time,
			ws.end_time,
			ws.grace_period_minutes
		FROM employees e
		INNER JOIN work_shifts ws ON ws.id = e.shift_id AND ws.deleted_at IS NULL
		WHERE e.tenant_id = ? AND e.user_id = ? AND e.deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"user_id":   userID,
		}).Msg("Failed to get work shift by user id")
		return nil, err
	}

	return &coreentity.WorkShift{
		ID:                 data.ID,
		TenantID:           data.TenantID,
		Name:               data.Name,
		Timezone:           data.Timezone,
		StartTime:          data.StartTime,
		EndTime:            data.EndTime,
		GracePeriodMinutes: data.GracePeriodMinutes,
	}, nil
}
