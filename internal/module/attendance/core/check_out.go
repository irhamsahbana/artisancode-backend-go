package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) CheckOut(ctx context.Context, data coreentity.AttendanceLogAction) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CheckOut")
	defer span.End()

	return c.createAttendanceLog(ctx, data, "check_out")
}
