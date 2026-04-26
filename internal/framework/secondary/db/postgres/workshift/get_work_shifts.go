package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *workShiftRepo) GetWorkShifts(
	ctx context.Context,
	filter coreentity.WorkShiftListFilter,
) ([]coreentity.WorkShift, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:workshift:get_work_shifts:GetWorkShifts",
	)
	defer span.End()

	type dao struct {
		TotalData          int    `db:"total_data"`
		ID                 string `db:"id"`
		TenantID           string `db:"tenant_id"`
		Name               string `db:"name"`
		Timezone           string `db:"timezone"`
		StartTime          string `db:"start_time"`
		EndTime            string `db:"end_time"`
		GracePeriodMinutes int    `db:"grace_period_minutes"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 4)
		items = make([]coreentity.WorkShift, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, timezone, start_time, end_time, grace_period_minutes
		FROM work_shifts
		WHERE deleted_at IS NULL AND tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query work shifts")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.WorkShift{
			ID:                 d.ID,
			TenantID:           d.TenantID,
			Name:               d.Name,
			Timezone:           d.Timezone,
			StartTime:          d.StartTime,
			EndTime:            d.EndTime,
			GracePeriodMinutes: d.GracePeriodMinutes,
		})
	}

	return items, total, nil
}
