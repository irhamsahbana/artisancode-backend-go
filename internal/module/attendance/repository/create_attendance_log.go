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
	ctx, span := tracing.StartSpan(ctx, "repo.CreateAttendanceLog")
	defer span.End()

	loggedAt, err := time.Parse(time.RFC3339, data.LoggedAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to parse attendance timestamp")
		return nil, err
	}

	query := `
		INSERT INTO attendance_logs (
			tenant_id, employee_id, attendance_date, type, source, status, logged_at,
			latitude, longitude, address, device_id, device_name, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	var (
		id        string
		createdAt time.Time
	)

	err = r.db.QueryRowxContext(
		ctx,
		r.db.Rebind(query),
		data.TenantID,
		data.EmployeeID,
		data.AttendanceDate,
		data.Type,
		data.Source,
		data.Status,
		loggedAt,
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
