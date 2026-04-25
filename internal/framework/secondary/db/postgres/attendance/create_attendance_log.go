package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) CreateAttendanceLog(ctx context.Context, data coreentity.AttendanceLog) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:attendance:create_attendance_log:CreateAttendanceLog")
	defer span.End()

	loggedAt, err := time.Parse(time.RFC3339, data.LoggedAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to parse attendance timestamp")
		return nil, err
	}

	query := `
		INSERT INTO attendance_logs (
			tenant_id, employee_id, attendance_date, type, source, status, logged_at,
			shift_id, shift_name, shift_timezone, shift_start_time, shift_end_time, shift_grace_period_minutes,
			latitude, longitude, address, device_id, device_name, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	var (
		id        string
		createdAt time.Time
	)

	exec := r.executor(ctx)
	err = exec.QueryRowxContext(
		ctx,
		exec.Rebind(query),
		data.TenantID,
		data.EmployeeID,
		data.AttendanceDate,
		string(data.Type),
		string(data.Source),
		string(data.Status),
		loggedAt,
		data.ShiftID,
		data.ShiftName,
		data.ShiftTimezone,
		data.ShiftStartTime,
		data.ShiftEndTime,
		data.ShiftGraceMins,
		data.Latitude,
		data.Longitude,
		data.Address,
		data.DeviceID,
		data.DeviceName,
		data.Notes,
	).Scan(&id, &createdAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create attendance log")
		return nil, err
	}

	data.ID = id
	data.CreatedAt = createdAt.Format(time.RFC3339)
	return &data, nil
}
