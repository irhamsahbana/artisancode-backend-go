package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *jobPositionRepo) UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:UpdateJobPosition")
	defer span.End()

	query := `
		UPDATE job_positions
		SET name = ?, grade = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.Name, data.Grade, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update job position")
		return err
	}
	return nil
}
