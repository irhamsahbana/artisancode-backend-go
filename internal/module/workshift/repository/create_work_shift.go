package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *workShiftRepo) CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateWorkShift")
	defer span.End()

	query := `
		INSERT INTO work_shifts (tenant_id, name, timezone, start_time, end_time, grace_period_minutes)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.Name, data.Timezone, data.StartTime, data.EndTime, data.GracePeriodMinutes,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create work shift")
		return nil, err
	}

	data.ID = id
	return &data, nil
}