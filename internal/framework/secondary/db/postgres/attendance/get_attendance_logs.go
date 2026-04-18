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

func (r *attendanceRepo) GetAttendanceLogs(ctx context.Context, filter coreentity.AttendanceLogListFilter) ([]coreentity.AttendanceLog, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:attendance:get_attendance_logs:GetAttendanceLogs")
	defer span.End()

	type dao struct {
		TotalData      int             `db:"total_data"`
		ID             string          `db:"id"`
		TenantID       string          `db:"tenant_id"`
		EmployeeID     string          `db:"employee_id"`
		EmployeeNo     string          `db:"employee_no"`
		EmployeeName   string          `db:"employee_name"`
		AttendanceDate time.Time       `db:"attendance_date"`
		ShiftID        *string         `db:"shift_id"`
		ShiftName      sql.NullString  `db:"shift_name"`
		ShiftTimezone  sql.NullString  `db:"shift_timezone"`
		ShiftStartTime sql.NullString  `db:"shift_start_time"`
		ShiftEndTime   sql.NullString  `db:"shift_end_time"`
		ShiftGraceMins sql.NullInt64   `db:"shift_grace_period_minutes"`
		Type           string          `db:"type"`
		Source         string          `db:"source"`
		Status         string          `db:"status"`
		LoggedAt       time.Time       `db:"logged_at"`
		Latitude       sql.NullFloat64 `db:"latitude"`
		Longitude      sql.NullFloat64 `db:"longitude"`
		Address        sql.NullString  `db:"address"`
		DeviceID       sql.NullString  `db:"device_id"`
		DeviceName     sql.NullString  `db:"device_name"`
		Notes          sql.NullString  `db:"notes"`
		CreatedAt      time.Time       `db:"created_at"`
		UpdatedAt      *time.Time      `db:"updated_at"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 10)
		items = make([]coreentity.AttendanceLog, 0)
		total = 0
	)

	query := `
		WITH daily_attendance AS (
			SELECT
				day_logs.employee_id,
				day_logs.attendance_date,
				MIN(CASE WHEN day_logs.type = 'check_in' THEN timezone(COALESCE(day_logs.shift_timezone, day_wl.timezone, 'UTC'), day_logs.logged_at) END) AS first_check_in_local,
				MAX(CASE WHEN day_logs.type = 'check_out' THEN timezone(COALESCE(day_logs.shift_timezone, day_wl.timezone, 'UTC'), day_logs.logged_at) END) AS last_check_out_local
			FROM attendance_logs day_logs
			INNER JOIN employees day_employee ON day_employee.id = day_logs.employee_id AND day_employee.deleted_at IS NULL
			LEFT JOIN work_locations day_wl ON day_wl.id = day_employee.location_id AND day_wl.deleted_at IS NULL
			WHERE day_logs.deleted_at IS NULL AND day_logs.tenant_id = ?
			GROUP BY day_logs.employee_id, day_logs.attendance_date
		)
		SELECT
			COUNT(*) OVER() AS total_data,
			al.id,
			al.tenant_id,
			al.employee_id,
			e.employee_no,
			e.full_name AS employee_name,
			al.attendance_date,
			al.shift_id,
			al.shift_name,
			al.shift_timezone,
			al.shift_start_time,
			al.shift_end_time,
			al.shift_grace_period_minutes,
			al.type,
			al.source,
			al.status,
			al.logged_at,
			al.latitude,
			al.longitude,
			al.address,
			al.device_id,
			al.device_name,
			al.notes,
			al.created_at,
			al.updated_at
		FROM attendance_logs al
		INNER JOIN employees e ON e.id = al.employee_id AND e.deleted_at IS NULL
		LEFT JOIN work_locations wl ON wl.id = e.location_id AND wl.deleted_at IS NULL
		LEFT JOIN daily_attendance da ON da.employee_id = al.employee_id AND da.attendance_date = al.attendance_date
		WHERE al.deleted_at IS NULL AND al.tenant_id = ?
	`
	args = append(args, filter.TenantID, filter.TenantID)

	if filter.EmployeeID != nil {
		query += ` AND al.employee_id = ?`
		args = append(args, *filter.EmployeeID)
	}
	if filter.Type != nil {
		query += ` AND al.type = ?`
		args = append(args, string(*filter.Type))
	}
	if filter.Source != nil {
		query += ` AND al.source = ?`
		args = append(args, string(*filter.Source))
	}
	if filter.Status != nil {
		query += ` AND al.status = ?`
		args = append(args, string(*filter.Status))
	}
	if filter.SelfieStatus != nil {
		if *filter.SelfieStatus == "with_photo" {
			query += ` AND EXISTS (
				SELECT 1
				FROM storage_file_links sfl
				WHERE sfl.tenant_id = al.tenant_id
					AND sfl.resource_type = 'attendance_log'
					AND sfl.resource_id = al.id
					AND sfl.field_name = 'selfie'
			)`
		}
		if *filter.SelfieStatus == "without_photo" {
			query += ` AND NOT EXISTS (
				SELECT 1
				FROM storage_file_links sfl
				WHERE sfl.tenant_id = al.tenant_id
					AND sfl.resource_type = 'attendance_log'
					AND sfl.resource_id = al.id
					AND sfl.field_name = 'selfie'
			)`
		}
	}
	if filter.OrgUnitID != nil {
		query += ` AND e.org_unit_id IS NOT NULL AND EXISTS (
			WITH RECURSIVE org_unit_ancestors AS (
				SELECT ou.id, ou.parent_id
				FROM org_units ou
				WHERE ou.id = e.org_unit_id AND ou.deleted_at IS NULL

				UNION ALL

				SELECT parent.id, parent.parent_id
				FROM org_units parent
				INNER JOIN org_unit_ancestors child ON child.parent_id = parent.id
				WHERE parent.deleted_at IS NULL
			)
			SELECT 1 FROM org_unit_ancestors WHERE id = ?
		)`
		args = append(args, *filter.OrgUnitID)
	}
	if filter.BranchID != nil {
		query += ` AND e.org_unit_id IS NOT NULL AND EXISTS (
			WITH RECURSIVE org_unit_ancestors AS (
				SELECT ou.id, ou.parent_id
				FROM org_units ou
				WHERE ou.id = e.org_unit_id AND ou.deleted_at IS NULL

				UNION ALL

				SELECT parent.id, parent.parent_id
				FROM org_units parent
				INNER JOIN org_unit_ancestors child ON child.parent_id = parent.id
				WHERE parent.deleted_at IS NULL
			)
			SELECT 1 FROM org_unit_ancestors WHERE id = ?
		)`
		args = append(args, *filter.BranchID)
	}
	if filter.WorkLocationID != nil {
		query += ` AND e.location_id = ?`
		args = append(args, *filter.WorkLocationID)
	}
	if filter.AttendanceDay != nil {
		query += ` AND al.attendance_date = ?`
		args = append(args, *filter.AttendanceDay)
	}
	if filter.DateFrom != nil {
		query += ` AND al.attendance_date >= ?`
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query += ` AND al.attendance_date <= ?`
		args = append(args, *filter.DateTo)
	}
	if filter.Q != "" {
		query += ` AND (e.full_name ILIKE '%' || ? || '%' OR e.employee_no ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	if filter.ExceptionType != nil {
		if *filter.ExceptionType == "late_check_in" {
			query += ` AND al.type = 'check_in'
				AND da.first_check_in_local IS NOT NULL
				AND al.shift_start_time IS NOT NULL
				AND al.shift_start_time <> ''
				AND da.first_check_in_local::time > (al.shift_start_time::time + make_interval(mins => COALESCE(al.shift_grace_period_minutes, 0)))
				AND timezone(COALESCE(al.shift_timezone, wl.timezone, 'UTC'), al.logged_at) = da.first_check_in_local`
		}
		if *filter.ExceptionType == "missing_check_out" {
			query += ` AND al.type = 'check_in'
				AND da.first_check_in_local IS NOT NULL
				AND da.last_check_out_local IS NULL
				AND timezone(COALESCE(al.shift_timezone, wl.timezone, 'UTC'), al.logged_at) = da.first_check_in_local`
		}
		if *filter.ExceptionType == "missing_check_in" {
			query += ` AND al.type = 'check_out'
				AND da.first_check_in_local IS NULL
				AND da.last_check_out_local IS NOT NULL
				AND timezone(COALESCE(al.shift_timezone, wl.timezone, 'UTC'), al.logged_at) = da.last_check_out_local`
		}
	}

	query += ` ORDER BY al.logged_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query attendance logs")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.AttendanceLog{
			ID:             d.ID,
			TenantID:       d.TenantID,
			EmployeeID:     d.EmployeeID,
			EmployeeNo:     d.EmployeeNo,
			EmployeeName:   d.EmployeeName,
			AttendanceDate: d.AttendanceDate.Format("2006-01-02"),
			ShiftID:        d.ShiftID,
			ShiftName:      nullableStringPtr(d.ShiftName),
			ShiftTimezone:  nullableStringPtr(d.ShiftTimezone),
			ShiftStartTime: nullableStringPtr(d.ShiftStartTime),
			ShiftEndTime:   nullableStringPtr(d.ShiftEndTime),
			ShiftGraceMins: nullableIntPtr(d.ShiftGraceMins),
			Type:           common.AttendanceType(d.Type),
			Source:         common.AttendanceSource(d.Source),
			Status:         common.AttendanceStatus(d.Status),
			LoggedAt:       d.LoggedAt.Format(time.RFC3339),
			Latitude:       nullableFloatToPtr(d.Latitude),
			Longitude:      nullableFloatToPtr(d.Longitude),
			Address:        nullableStringPtr(d.Address),
			DeviceID:       nullableStringPtr(d.DeviceID),
			DeviceName:     nullableStringPtr(d.DeviceName),
			Notes:          nullableStringPtr(d.Notes),
			CreatedAt:      d.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      formatTimePtr(d.UpdatedAt),
		})
	}

	return items, total, nil
}

func nullableFloatToPtr(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}

	result := value.Float64
	return &result
}

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	result := value.String
	return &result
}

func formatTimePtr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.Format(time.RFC3339)
	return &formatted
}

func nullableIntPtr(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}

	result := int(value.Int64)
	return &result
}
