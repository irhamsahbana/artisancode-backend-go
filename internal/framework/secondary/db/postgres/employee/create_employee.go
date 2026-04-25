package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:create_employee:CreateEmployee")
	defer span.End()

	query := `
		INSERT INTO employees (
			tenant_id, employee_no, full_name, user_id, org_unit_id, job_position_id,
			location_id, shift_id, email, status, join_date, join_date_timezone
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &id, exec.Rebind(query),
		data.TenantID,
		data.EmployeeNo,
		data.FullName,
		data.UserID,
		data.OrgUnitID,
		data.JobPositionID,
		data.LocationID,
		data.ShiftID,
		data.Email,
		data.Status,
		data.JoinDate,
		data.JoinDateTimezone,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create employee")
		return nil, err
	}

	data.ID = id
	return &data, nil
}
