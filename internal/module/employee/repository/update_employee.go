package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) UpdateEmployee(ctx context.Context, data coreentity.Employee) error {
	ctx, span := tracing.StartSpan(ctx, "repo.UpdateEmployee")
	defer span.End()

	query := `
		UPDATE employees
		SET employee_no = ?, full_name = ?, user_id = ?, org_unit_id = ?,
			job_position_id = ?, location_id = ?, shift_id = ?,
			status = ?, join_date = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.EmployeeNo,
		data.FullName,
		data.UserID,
		data.OrgUnitID,
		data.JobPositionID,
		data.LocationID,
		data.ShiftID,
		data.Status,
		data.JoinDate,
		data.ID,
		data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update employee")
		return err
	}
	return nil
}