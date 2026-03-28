package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo.DeleteEmployee")
	defer span.End()

	query := `
		UPDATE employees
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete employee")
		return err
	}
	return nil
}