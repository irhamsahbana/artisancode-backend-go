package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetAttendanceLog(ctx context.Context, filter coreentity.AttendanceLogDetailFilter) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:attendance:get_attendance_log:GetAttendanceLog")
	defer span.End()

	var data struct {
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

	query := `
		SELECT
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
		WHERE al.id = ? AND al.tenant_id = ? AND al.deleted_at IS NULL
	`

	args := []any{filter.ID, filter.TenantID}
	if filter.EmployeeID != nil {
		query += ` AND al.employee_id = ?`
		args = append(args, *filter.EmployeeID)
	}

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Attendance log not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get attendance log")
		return nil, err
	}

	return &coreentity.AttendanceLog{
		ID:             data.ID,
		TenantID:       data.TenantID,
		EmployeeID:     data.EmployeeID,
		EmployeeNo:     data.EmployeeNo,
		EmployeeName:   data.EmployeeName,
		AttendanceDate: data.AttendanceDate.Format("2006-01-02"),
		Type:           common.AttendanceType(data.Type),
		Source:         common.AttendanceSource(data.Source),
		Status:         common.AttendanceStatus(data.Status),
		LoggedAt:       data.LoggedAt.Format(time.RFC3339),
		Latitude:       nullableFloatToPtr(data.Latitude),
		Longitude:      nullableFloatToPtr(data.Longitude),
		Address:        nullableStringPtr(data.Address),
		DeviceID:       nullableStringPtr(data.DeviceID),
		DeviceName:     nullableStringPtr(data.DeviceName),
		Notes:          nullableStringPtr(data.Notes),
		CreatedAt:      data.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      formatTimePtr(data.UpdatedAt),
	}, nil
}
