package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) GetAttendanceLogs(
	ctx context.Context,
	filter coreentity.AttendanceLogListFilter,
) ([]coreentity.AttendanceLog, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:attendance:get_attendance_logs:GetAttendanceLogs")
	defer span.End()

	if !filter.UserCtx.HasRole("owner") {
		employee, err := c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserCtx.UserID)
		if err != nil {
			return []coreentity.AttendanceLog{}, 0, nil
		}
		filter.EmployeeID = &employee.ID
	}

	items, total, err := c.repo.GetAttendanceLogs(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for index := range items {
		err = c.attachAttendanceLogSelfie(ctx, &items[index])
		if err != nil {
			return nil, 0, err
		}
	}

	return items, total, nil
}
