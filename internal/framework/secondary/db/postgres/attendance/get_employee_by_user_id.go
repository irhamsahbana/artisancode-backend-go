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

func (r *attendanceRepo) GetEmployeeByUserID(
	ctx context.Context,
	tenantID, userID string,
) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:attendance:get_employee_by_user_id:GetEmployeeByUserID",
	)
	defer span.End()

	var data struct {
		ID         string         `db:"id"`
		TenantID   string         `db:"tenant_id"`
		EmployeeNo string         `db:"employee_no"`
		FullName   string         `db:"full_name"`
		Email      sql.NullString `db:"email"`
		LocationID *string        `db:"location_id"`
		ShiftID    *string        `db:"shift_id"`
	}

	query := `
		SELECT id, tenant_id, employee_no, full_name, email, location_id, shift_id
		FROM employees
		WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Employee profile not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"user_id":   userID,
		}).Msg("Failed to get employee by user id")
		return nil, err
	}

	return &coreentity.Employee{
		ID:         data.ID,
		TenantID:   data.TenantID,
		EmployeeNo: data.EmployeeNo,
		FullName:   data.FullName,
		Email:      nullableStringToValue(data.Email),
		LocationID: data.LocationID,
		ShiftID:    data.ShiftID,
	}, nil
}

func nullableStringToValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}

	return value.String
}
