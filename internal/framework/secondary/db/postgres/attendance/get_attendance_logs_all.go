package repository

import (
	"codebase-app/internal/infrastructure/tracing"
	"context"

	"codebase-app/internal/entity/coreentity"
)

func (r *attendanceRepo) GetAttendanceLogsAll(
	ctx context.Context,
	filter coreentity.AttendanceLogListFilter,
) ([]coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:attendance:get_attendance_logs_all:GetAttendanceLogsAll",
	)
	defer span.End()

	filter.Page = 1
	filter.Paginate = 1000000

	items, _, err := r.GetAttendanceLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	return items, nil
}
