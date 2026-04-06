package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetOwnerAttendanceDashboardTrend(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) ([]coreentity.OwnerAttendanceDashboardDailyTrend, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetOwnerAttendanceDashboardTrend")
	defer span.End()

	type dao struct {
		AttendanceDate       time.Time `db:"attendance_date"`
		CheckedInCount       int       `db:"checked_in_count"`
		CheckedOutCount      int       `db:"checked_out_count"`
		LateCheckInCount     int       `db:"late_check_in_count"`
		MissingCheckOutCount int       `db:"missing_check_out_count"`
	}

	query := `
		WITH active_employees AS (
			SELECT
				e.id,
				ws.start_time,
				COALESCE(ws.grace_period_minutes, 0) AS grace_period_minutes
			FROM employees e
			LEFT JOIN work_shifts ws ON ws.id = e.shift_id AND ws.deleted_at IS NULL
			WHERE e.tenant_id = ? AND e.deleted_at IS NULL AND e.status = 'active'
		),
		date_series AS (
			SELECT generate_series((?::date - (?::int - 1) * INTERVAL '1 day')::date, ?::date, INTERVAL '1 day')::date AS attendance_date
		),
		daily_logs AS (
			SELECT
				al.attendance_date,
				al.employee_id,
				MIN(CASE WHEN al.type = 'check_in' THEN al.logged_at END) AS first_check_in_at,
				MAX(CASE WHEN al.type = 'check_out' THEN al.logged_at END) AS last_check_out_at
			FROM attendance_logs al
			INNER JOIN active_employees ae ON ae.id = al.employee_id
			WHERE
				al.tenant_id = ?
				AND al.deleted_at IS NULL
				AND al.attendance_date BETWEEN (?::date - (?::int - 1) * INTERVAL '1 day')::date AND ?::date
			GROUP BY al.attendance_date, al.employee_id
		)
		SELECT
			ds.attendance_date,
			COUNT(*) FILTER (WHERE dl.first_check_in_at IS NOT NULL) AS checked_in_count,
			COUNT(*) FILTER (WHERE dl.last_check_out_at IS NOT NULL) AS checked_out_count,
			COUNT(*) FILTER (
				WHERE dl.first_check_in_at IS NOT NULL
					AND ae.start_time IS NOT NULL
					AND ae.start_time <> ''
					AND dl.first_check_in_at::time > (ae.start_time::time + make_interval(mins => ae.grace_period_minutes))
			) AS late_check_in_count,
			COUNT(*) FILTER (WHERE dl.first_check_in_at IS NOT NULL AND dl.last_check_out_at IS NULL) AS missing_check_out_count
		FROM date_series ds
		CROSS JOIN active_employees ae
		LEFT JOIN daily_logs dl ON dl.attendance_date = ds.attendance_date AND dl.employee_id = ae.id
		GROUP BY ds.attendance_date
		ORDER BY ds.attendance_date ASC
	`

	rows := make([]dao, 0)
	err := r.db.SelectContext(
		ctx,
		&rows,
		r.db.Rebind(query),
		filter.TenantID,
		filter.Date,
		filter.TrendDays,
		filter.Date,
		filter.TenantID,
		filter.Date,
		filter.TrendDays,
		filter.Date,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query owner attendance dashboard trend")
		return nil, err
	}

	items := make([]coreentity.OwnerAttendanceDashboardDailyTrend, 0, len(rows))
	for _, row := range rows {
		items = append(items, coreentity.OwnerAttendanceDashboardDailyTrend{
			UserCtx:              filter.UserCtx,
			AttendanceDate:       row.AttendanceDate.Format("2006-01-02"),
			CheckedInCount:       row.CheckedInCount,
			CheckedOutCount:      row.CheckedOutCount,
			LateCheckInCount:     row.LateCheckInCount,
			MissingCheckOutCount: row.MissingCheckOutCount,
		})
	}

	return items, nil
}
