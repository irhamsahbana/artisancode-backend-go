package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetOwnerAttendanceDashboardSummary(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) (*coreentity.OwnerAttendanceDashboardSummary, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:attendance:get_owner_attendance_dashboard_summary:GetOwnerAttendanceDashboardSummary")
	defer span.End()

	var data struct {
		ActiveEmployeeCount  int `db:"active_employee_count"`
		CheckedInCount       int `db:"checked_in_count"`
		CheckedOutCount      int `db:"checked_out_count"`
		PendingCheckInCount  int `db:"pending_check_in_count"`
		PendingCheckOutCount int `db:"pending_check_out_count"`
		LateCheckInCount     int `db:"late_check_in_count"`
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
			COUNT(*) AS active_employee_count,
			COUNT(*) FILTER (WHERE dl.first_check_in_at_local IS NOT NULL) AS checked_in_count,
			COUNT(*) FILTER (WHERE dl.last_check_out_at_local IS NOT NULL) AS checked_out_count,
			COUNT(*) FILTER (WHERE dl.first_check_in_at_local IS NULL) AS pending_check_in_count,
			COUNT(*) FILTER (WHERE dl.first_check_in_at_local IS NOT NULL AND dl.last_check_out_at_local IS NULL) AS pending_check_out_count,
			COUNT(*) FILTER (
				WHERE dl.first_check_in_at_local IS NOT NULL
					AND ae.start_time IS NOT NULL
					AND ae.start_time <> ''
					AND dl.first_check_in_at_local::time > (ae.start_time::time + make_interval(mins => ae.grace_period_minutes))
			) AS late_check_in_count
		FROM active_employees ae
		LEFT JOIN daily_logs dl ON dl.employee_id = ae.id
	`

	err := r.db.GetContext(
		ctx,
		&data,
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
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query owner attendance dashboard summary")
		return nil, err
	}

	return &coreentity.OwnerAttendanceDashboardSummary{
		UserCtx:              filter.UserCtx,
		ActiveEmployeeCount:  data.ActiveEmployeeCount,
		CheckedInCount:       data.CheckedInCount,
		CheckedOutCount:      data.CheckedOutCount,
		PendingCheckInCount:  data.PendingCheckInCount,
		PendingCheckOutCount: data.PendingCheckOutCount,
		LateCheckInCount:     data.LateCheckInCount,
	}, nil
}
