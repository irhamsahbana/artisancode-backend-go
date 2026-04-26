package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:attendance:get_work_shift:GetWorkShift",
	)
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
		SELECT id, tenant_id, name, timezone, start_time, end_time, grace_period_minutes
		FROM work_shifts
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Work shift not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get work shift")
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
