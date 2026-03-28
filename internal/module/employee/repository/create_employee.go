package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateEmployee")
	defer span.End()

	query := `
		INSERT INTO employees (
			tenant_id, employee_no, full_name, user_id, org_unit_id, job_position_id,
			location_id, shift_id, status, join_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID,
		data.EmployeeNo,
		data.FullName,
		data.UserID,
		data.OrgUnitID,
		data.JobPositionID,
		data.LocationID,
		data.ShiftID,
		data.Status,
		data.JoinDate,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create employee")
		return nil, err
	}

	data.ID = id
	return &data, nil
}