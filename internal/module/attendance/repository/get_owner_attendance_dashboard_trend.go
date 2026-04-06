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
		WITH trend_window AS (
			SELECT
				?::date AS selected_date,
				?::text AS timezone_name,
				(?::int - 1) AS days_back,
				((?::date - (?::int - 1) * INTERVAL '1 day')::timestamp AT TIME ZONE ?::text) AS utc_window_start,
				((?::date + INTERVAL '1 day')::timestamp AT TIME ZONE ?::text) AS utc_window_end
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
		date_series AS (
			SELECT generate_series((tw.selected_date - tw.days_back * INTERVAL '1 day')::date, tw.selected_date, INTERVAL '1 day')::date AS attendance_date
			FROM trend_window tw
		),
		daily_logs AS (
			SELECT
				timezone(tw.timezone_name, al.logged_at)::date AS attendance_date_local,
				al.employee_id,
				MIN(CASE WHEN al.type = 'check_in' THEN timezone(tw.timezone_name, al.logged_at) END) AS first_check_in_at_local,
				MAX(CASE WHEN al.type = 'check_out' THEN timezone(tw.timezone_name, al.logged_at) END) AS last_check_out_at_local
			FROM attendance_logs al
			INNER JOIN active_employees ae ON ae.id = al.employee_id
			CROSS JOIN trend_window tw
			WHERE
				al.tenant_id = ?
				AND al.deleted_at IS NULL
				AND al.logged_at >= tw.utc_window_start
				AND al.logged_at < tw.utc_window_end
			GROUP BY timezone(tw.timezone_name, al.logged_at)::date, al.employee_id
		)
		SELECT
			ds.attendance_date,
			COUNT(*) FILTER (WHERE dl.first_check_in_at_local IS NOT NULL) AS checked_in_count,
			COUNT(*) FILTER (WHERE dl.last_check_out_at_local IS NOT NULL) AS checked_out_count,
			COUNT(*) FILTER (
				WHERE dl.first_check_in_at_local IS NOT NULL
					AND ae.start_time IS NOT NULL
					AND ae.start_time <> ''
					AND dl.first_check_in_at_local::time > (ae.start_time::time + make_interval(mins => ae.grace_period_minutes))
			) AS late_check_in_count,
			COUNT(*) FILTER (WHERE dl.first_check_in_at_local IS NOT NULL AND dl.last_check_out_at_local IS NULL) AS missing_check_out_count
		FROM date_series ds
		CROSS JOIN active_employees ae
		LEFT JOIN daily_logs dl ON dl.attendance_date_local = ds.attendance_date AND dl.employee_id = ae.id
		GROUP BY ds.attendance_date
		ORDER BY ds.attendance_date ASC
	`

	rows := make([]dao, 0)
	err := r.db.SelectContext(
		ctx,
		&rows,
		r.db.Rebind(query),
		filter.Date,
		filter.Timezone,
		filter.TrendDays,
		filter.Date,
		filter.TrendDays,
		filter.Timezone,
		filter.Date,
		filter.Timezone,
		filter.TenantID,
		filter.TenantID,
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
