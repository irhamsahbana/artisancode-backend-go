package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) CheckOut(ctx context.Context, data coreentity.AttendanceLogAction) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:attendance:check_out:CheckOut")
	defer span.End()

	return c.createAttendanceLog(ctx, data, common.AttendanceTypeCheckOut)
}
