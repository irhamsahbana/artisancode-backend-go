package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type AttendanceRepository interface {
	GetAttendanceLogs(ctx context.Context, filter coreentity.AttendanceLogListFilter) ([]coreentity.AttendanceLog, int, error)
	GetAttendanceLogsAll(ctx context.Context, filter coreentity.AttendanceLogListFilter) ([]coreentity.AttendanceLog, error)
	GetAttendanceLog(ctx context.Context, filter coreentity.AttendanceLogDetailFilter) (*coreentity.AttendanceLog, error)
	CreateAttendanceLog(ctx context.Context, data coreentity.AttendanceLog) (*coreentity.AttendanceLog, error)
	ExistsAttendanceByTypeOnDate(ctx context.Context, tenantID, employeeID, attendanceDate, attendanceType string) (bool, error)
	GetEmployeeByUserID(ctx context.Context, tenantID, userID string) (*coreentity.Employee, error)
	GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error)
	GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	GetOwnerAttendanceDashboardSummary(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) (*coreentity.OwnerAttendanceDashboardSummary, error)
	GetOwnerAttendanceDashboardExceptions(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) ([]coreentity.OwnerAttendanceDashboardException, error)
	GetOwnerAttendanceDashboardTrend(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) ([]coreentity.OwnerAttendanceDashboardDailyTrend, error)
}
