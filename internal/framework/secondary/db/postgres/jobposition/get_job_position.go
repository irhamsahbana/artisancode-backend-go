package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *jobPositionRepo) GetJobPosition(
	ctx context.Context,
	filter coreentity.JobPosition,
) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:GetJobPosition")
	defer span.End()

	var data = new(repoentity.JobPosition)

	query := `
		SELECT id, tenant_id, name, grade
		FROM job_positions
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageJobPositionNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageJobPositionNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get job position")
		return nil, err
	}

	result := &coreentity.JobPosition{
		ID:       data.ID,
		TenantID: data.TenantID,
		Name:     data.Name,
		Grade:    data.Grade,
	}
	return result, nil
}
