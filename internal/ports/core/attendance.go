package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type AttendanceCore interface {
	GetAttendanceLogs(ctx context.Context, filter coreentity.AttendanceLogListFilter) ([]coreentity.AttendanceLog, int, error)
	GetAttendanceLog(ctx context.Context, filter coreentity.AttendanceLogDetailFilter) (*coreentity.AttendanceLog, error)
	CheckIn(ctx context.Context, data coreentity.AttendanceLogAction) (*coreentity.AttendanceLog, error)
	CheckOut(ctx context.Context, data coreentity.AttendanceLogAction) (*coreentity.AttendanceLog, error)
	GetAttendanceSummaryToday(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.AttendanceSummary, error)
	GetAttendancePolicy(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.AttendancePolicy, error)
}
