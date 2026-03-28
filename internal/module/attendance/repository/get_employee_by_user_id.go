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

func (r *attendanceRepo) GetEmployeeByUserID(ctx context.Context, tenantID, userID string) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetEmployeeByUserID")
	defer span.End()

	var data struct {
		ID         string         `db:"id"`
		TenantID   string         `db:"tenant_id"`
		EmployeeNo string         `db:"employee_no"`
		FullName   string         `db:"full_name"`
		Email      sql.NullString `db:"email"`
	}

	query := `
		SELECT id, tenant_id, employee_no, full_name, email
		FROM employees
		WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), tenantID, userID)
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
	}, nil
}

func nullableStringToValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}

	return value.String
}
