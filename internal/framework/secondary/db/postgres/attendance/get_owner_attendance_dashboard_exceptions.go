package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetOwnerAttendanceDashboardExceptions(
	ctx context.Context,
	filter coreentity.OwnerAttendanceDashboardFilter,
) ([]coreentity.OwnerAttendanceDashboardException, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:attendance:get_owner_attendance_dashboard_exceptions:GetOwnerAttendanceDashboardExceptions",
	)
	defer span.End()

	type dao struct {
		EmployeeID     string         `db:"employee_id"`
		EmployeeNo     string         `db:"employee_no"`
		EmployeeName   string         `db:"employee_name"`
		ShiftName      sql.NullString `db:"shift_name"`
		FirstCheckInAt *time.Time     `db:"first_check_in_at"`
		LastCheckOutAt *time.Time     `db:"last_check_out_at"`
		ExceptionType  string         `db:"exception_type"`
	}

	query := `
		WITH dashboard_day AS (
			SELECT
				?::date AS selected_date,
				?::text AS timezone_name,
				(?::date::timestamp AT TIME ZONE ?::text) AS utc_start,
				((?::date + INTERVAL '1 day')::timestamp AT TIME ZONE ?::text) AS utc_end
		),
		active_employees AS (
			SELECT
				e.id,
				e.employee_no,
				e.full_name,
				ws.name AS shift_name,
				ws.start_time,
				COALESCE(ws.grace_period_minutes, 0) AS grace_period_minutes
			FROM employees e
			LEFT JOIN work_shifts ws ON ws.id = e.shift_id AND ws.deleted_at IS NULL
			WHERE e.tenant_id = ? AND e.deleted_at IS NULL AND e.status = 'active'
		),
		daily_logs AS (
			SELECT
				al.employee_id,
				MIN(CASE WHEN al.type = 'check_in' THEN timezone(dd.timezone_name, al.logged_at) END) AS first_check_in_at_local,
				MAX(CASE WHEN al.type = 'check_out' THEN timezone(dd.timezone_name, al.logged_at) END) AS last_check_out_at_local
			FROM attendance_logs al
			INNER JOIN active_employees ae ON ae.id = al.employee_id
			CROSS JOIN dashboard_day dd
			WHERE
				al.tenant_id = ?
				AND al.deleted_at IS NULL
				AND al.logged_at >= dd.utc_start
				AND al.logged_at < dd.utc_end
			GROUP BY al.employee_id
		)
		SELECT
			ae.id AS employee_id,
			ae.employee_no,
			ae.full_name AS employee_name,
			ae.shift_name,
			dl.first_check_in_at_local AS first_check_in_at,
			dl.last_check_out_at_local AS last_check_out_at,
			CASE
				WHEN dl.first_check_in_at_local IS NULL THEN 'missing_check_in'
				WHEN dl.last_check_out_at_local IS NULL THEN 'missing_check_out'
				WHEN ae.start_time IS NOT NULL
					AND ae.start_time <> ''
					AND dl.first_check_in_at_local::time > (ae.start_time::time + make_interval(mins => ae.grace_period_minutes))
				THEN 'late_check_in'
				ELSE ''
			END AS exception_type
		FROM active_employees ae
		LEFT JOIN daily_logs dl ON dl.employee_id = ae.id
		WHERE
			dl.first_check_in_at_local IS NULL
			OR (dl.first_check_in_at_local IS NOT NULL AND dl.last_check_out_at_local IS NULL)
			OR (
				dl.first_check_in_at_local IS NOT NULL
				AND ae.start_time IS NOT NULL
				AND ae.start_time <> ''
				AND dl.first_check_in_at_local::time > (ae.start_time::time + make_interval(mins => ae.grace_period_minutes))
			)
		ORDER BY
			CASE
				WHEN dl.first_check_in_at_local IS NULL THEN 1
				WHEN dl.last_check_out_at_local IS NULL THEN 2
				ELSE 3
			END,
			ae.full_name ASC
	`

	rows := make([]dao, 0)
	err := r.db.SelectContext(
		ctx,
		&rows,
		r.db.Rebind(query),
		filter.Date,
		filter.Timezone,
		filter.Date,
		filter.Timezone,
		filter.Date,
		filter.Timezone,
		filter.TenantID,
		filter.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, filter).
			Msg("Failed to query owner attendance dashboard exceptions")
		return nil, err
	}

	items := make([]coreentity.OwnerAttendanceDashboardException, 0, len(rows))
	for _, row := range rows {
		items = append(items, coreentity.OwnerAttendanceDashboardException{
			UserCtx:        filter.UserCtx,
			EmployeeID:     row.EmployeeID,
			EmployeeNo:     row.EmployeeNo,
			EmployeeName:   row.EmployeeName,
			ShiftName:      nullableStringPtr(row.ShiftName),
			FirstCheckInAt: formatTimePtr(row.FirstCheckInAt),
			LastCheckOutAt: formatTimePtr(row.LastCheckOutAt),
			ExceptionType:  row.ExceptionType,
		})
	}

	return items, nil
}
