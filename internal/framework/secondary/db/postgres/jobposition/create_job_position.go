package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *jobPositionRepo) CreateJobPosition(
	ctx context.Context,
	data coreentity.JobPosition,
) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:CreateJobPosition")
	defer span.End()

	query := `
		INSERT INTO job_positions (tenant_id, name, grade)
		VALUES (?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query), data.TenantID, data.Name, data.Grade)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create job position")
		return nil, err
	}

	data.ID = id
	return &data, nil
}
