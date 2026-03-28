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
	ctx, span := tracing.StartSpan(ctx, "repo.GetEmployee")
	defer span.End()

	var data struct {
		ID            string  `db:"id"`
		TenantID      string  `db:"tenant_id"`
		EmployeeNo    string  `db:"employee_no"`
		FullName      string  `db:"full_name"`
		Email         sql.NullString `db:"email"`
		UserID        *string `db:"user_id"`
		OrgUnitID     *string `db:"org_unit_id"`
		JobPositionID *string `db:"job_position_id"`
		LocationID    *string `db:"location_id"`
		ShiftID       *string `db:"shift_id"`
		Status        string  `db:"status"`
		JoinDate      *string `db:"join_date"`
	}

	query := `
		SELECT id, tenant_id, employee_no, full_name, email, user_id, org_unit_id, job_position_id,
			location_id, shift_id, status, join_date
		FROM employees
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Employee not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get employee")
		return nil, err
	}

	return &coreentity.Employee{
		ID:            data.ID,
		TenantID:      data.TenantID,
		EmployeeNo:    data.EmployeeNo,
		FullName:      data.FullName,
		Email:         nullableStringToValue(data.Email),
		UserID:        data.UserID,
		OrgUnitID:     data.OrgUnitID,
		JobPositionID: data.JobPositionID,
		LocationID:    data.LocationID,
		ShiftID:       data.ShiftID,
		Status:        data.Status,
		JoinDate:      data.JoinDate,
	}, nil
}
