package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *employeeRepo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetEmployees")
	defer span.End()

	type dao struct {
		TotalData        int            `db:"total_data"`
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
		JoinDate         *string        `db:"join_date"`
		JoinDateTimezone *string        `db:"join_date_timezone"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 8)
		items = make([]coreentity.Employee, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, employee_no, full_name, email, user_id, org_unit_id, job_position_id,
			location_id, shift_id, status, join_date, join_date_timezone
		FROM employees
		WHERE deleted_at IS NULL AND tenant_id = ?
	`
	args = append(args, filter.TenantID)

	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	if filter.OrgUnitID != nil {
		query += ` AND org_unit_id = ?`
		args = append(args, filter.OrgUnitID)
	}
	if filter.Q != "" {
		query += ` AND (full_name ILIKE '%' || ? || '%' OR employee_no ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}

	query += ` ORDER BY full_name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query employees")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.Employee{
			ID:               d.ID,
			TenantID:         d.TenantID,
			EmployeeNo:       d.EmployeeNo,
			FullName:         d.FullName,
			Email:            nullableStringToValue(d.Email),
			UserID:           d.UserID,
			OrgUnitID:        d.OrgUnitID,
			JobPositionID:    d.JobPositionID,
			LocationID:       d.LocationID,
			ShiftID:          d.ShiftID,
			Status:           d.Status,
			JoinDate:         d.JoinDate,
			JoinDateTimezone: d.JoinDateTimezone,
		})
	}

	return items, total, nil
}

func nullableStringToValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}

	return value.String
}
