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

func (r *employeeRepo) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:get_employee:GetEmployee")
	defer span.End()

	var data struct {
		ID               string         `db:"id"`
		TenantID         string         `db:"tenant_id"`
		EmployeeNo       string         `db:"employee_no"`
		FullName         string         `db:"full_name"`
		Email            sql.NullString `db:"email"`
		UserID           *string        `db:"user_id"`
		OrgUnitID        *string        `db:"org_unit_id"`
		JobPositionID    *string        `db:"job_position_id"`
		LocationID       *string        `db:"location_id"`
		ShiftID          *string        `db:"shift_id"`
		Status           string         `db:"status"`
		AccessStatus     string         `db:"access_status"`
		JoinDate         *string        `db:"join_date"`
		JoinDateTimezone *string        `db:"join_date_timezone"`
	}

	query := `
		SELECT id, tenant_id, employee_no, full_name, email, user_id, org_unit_id, job_position_id,
			location_id, shift_id, status,
			CASE
				WHEN user_id IS NOT NULL THEN 'active'
				WHEN EXISTS (
					SELECT 1
					FROM user_invitations ui
					WHERE ui.employee_id = employees.id
						AND ui.status = 'pending'
						AND ui.expires_at >= NOW()
						AND ui.deleted_at IS NULL
				) THEN 'invited'
				ELSE 'no_access'
			END AS access_status,
			join_date, join_date_timezone
		FROM employees
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Employee not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get employee")
		return nil, err
	}

	return &coreentity.Employee{
		ID:               data.ID,
		TenantID:         data.TenantID,
		EmployeeNo:       data.EmployeeNo,
		FullName:         data.FullName,
		Email:            nullableStringToValue(data.Email),
		UserID:           data.UserID,
		OrgUnitID:        data.OrgUnitID,
		JobPositionID:    data.JobPositionID,
		LocationID:       data.LocationID,
		ShiftID:          data.ShiftID,
		Status:           data.Status,
		AccessStatus:     data.AccessStatus,
		JoinDate:         data.JoinDate,
		JoinDateTimezone: data.JoinDateTimezone,
	}, nil
}
