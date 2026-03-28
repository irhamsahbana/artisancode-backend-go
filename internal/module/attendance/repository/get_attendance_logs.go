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
	ctx, span := tracing.StartSpan(ctx, "repo.GetAttendanceLogs")
	defer span.End()

	type dao struct {
		TotalData      int             `db:"total_data"`
		ID             string          `db:"id"`
		TenantID       string          `db:"tenant_id"`
		EmployeeID     string          `db:"employee_id"`
		EmployeeNo     string          `db:"employee_no"`
		EmployeeName   string          `db:"employee_name"`
		AttendanceDate time.Time       `db:"attendance_date"`
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
		SELECT
			COUNT(*) OVER() AS total_data,
			al.id,
			al.tenant_id,
			al.employee_id,
			e.employee_no,
			e.full_name AS employee_name,
			al.attendance_date,
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
		WHERE al.deleted_at IS NULL AND al.tenant_id = ?
	`
	args = append(args, filter.TenantID)

	if filter.EmployeeID != nil {
		query += ` AND al.employee_id = ?`
		args = append(args, *filter.EmployeeID)
	}
	if filter.Type != "" {
		query += ` AND al.type = ?`
		args = append(args, filter.Type)
	}
	if filter.Source != "" {
		query += ` AND al.source = ?`
		args = append(args, filter.Source)
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
			Type:           d.Type,
			Source:         d.Source,
			Status:         d.Status,
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
