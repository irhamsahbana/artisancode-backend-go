package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) GetAttendanceLog(ctx context.Context, filter coreentity.AttendanceLogDetailFilter) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetAttendanceLog")
	defer span.End()

	if !filter.UserCtx.IsOwner() {
		employee, err := c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserCtx.UserID)
		if err != nil {
			return nil, err
		}
		filter.EmployeeID = &employee.ID
	}

	return c.repo.GetAttendanceLog(ctx, filter)
}
