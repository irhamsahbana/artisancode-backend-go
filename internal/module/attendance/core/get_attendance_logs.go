package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) GetAttendanceLogs(ctx context.Context, filter coreentity.AttendanceLogListFilter) ([]coreentity.AttendanceLog, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetAttendanceLogs")
	defer span.End()

	if !filter.UserCtx.IsOwner() {
		employee, err := c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserCtx.UserID)
		if err != nil {
			return []coreentity.AttendanceLog{}, 0, nil
		}
		filter.EmployeeID = &employee.ID
	}

	return c.repo.GetAttendanceLogs(ctx, filter)
}
