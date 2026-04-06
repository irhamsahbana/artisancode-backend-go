package coreentity

import "codebase-app/internal/entity/common"

type AttendanceLog struct {
	UserCtx common.UserContext

	ID             string
	TenantID       string
	EmployeeID     string
	EmployeeNo     string
	EmployeeName   string
	AttendanceDate string
	Type           common.AttendanceType
	Source         common.AttendanceSource
	Status         common.AttendanceStatus
	LoggedAt       string
	Latitude       *float64
	Longitude      *float64
	Address        *string
	DeviceID       *string
	DeviceName     *string
	Notes          *string
	SelfieFileID   *string
	SelfieURL      *string
	CreatedAt      string
	UpdatedAt      *string
}

type AttendanceLogListFilter struct {
	UserCtx       common.UserContext
	TenantID      string
	EmployeeID    *string
	Q             string
	Type          *common.AttendanceType
	Source        *common.AttendanceSource
	AttendanceDay *string
	DateFrom      *string
	DateTo        *string
	Page          int
	Paginate      int
}

type AttendanceLogDetailFilter struct {
	UserCtx    common.UserContext
	TenantID   string
	ID         string
	EmployeeID *string
}

type AttendanceLogAction struct {
	UserCtx      common.UserContext
	TenantID     string
	Type         common.AttendanceType
	LoggedAt     string
	Latitude     *float64
	Longitude    *float64
	Address      *string
	DeviceID     *string
	DeviceName   *string
	Notes        *string
	SelfieFileID string
}

type AttendanceSummary struct {
	UserCtx common.UserContext

	AttendanceDate string
	CheckedIn      bool
	CheckedOut     bool
	CheckInLogID   *string
	CheckOutLogID  *string
	LastLogType    *common.AttendanceType
	LastLoggedAt   *string
	CanCheckIn     bool
	CanCheckOut    bool
}

type AttendancePolicy struct {
	UserCtx common.UserContext

	Timezone                string
	AttendanceRadiusMeters  int
	AttendanceCheckInStart  string
	AttendanceCheckInEnd    string
	AttendanceCheckOutStart string
	AttendanceCheckOutEnd   string
}

type OwnerAttendanceDashboardFilter struct {
	UserCtx common.UserContext

	TenantID  string
	Date      string
	Timezone  string
	TrendDays int
}

type OwnerAttendanceDashboard struct {
	UserCtx common.UserContext

	AttendanceDate  string
	Summary         OwnerAttendanceDashboardSummary
	TodayExceptions []OwnerAttendanceDashboardException
	DailyTrend      []OwnerAttendanceDashboardDailyTrend
}

type OwnerAttendanceDashboardSummary struct {
	UserCtx common.UserContext

	ActiveEmployeeCount  int
	CheckedInCount       int
	CheckedOutCount      int
	PendingCheckInCount  int
	PendingCheckOutCount int
	LateCheckInCount     int
}

type OwnerAttendanceDashboardException struct {
	UserCtx common.UserContext

	EmployeeID     string
	EmployeeNo     string
	EmployeeName   string
	ShiftName      *string
	FirstCheckInAt *string
	LastCheckOutAt *string
	ExceptionType  string
}

type OwnerAttendanceDashboardDailyTrend struct {
	UserCtx common.UserContext

	AttendanceDate       string
	CheckedInCount       int
	CheckedOutCount      int
	LateCheckInCount     int
	MissingCheckOutCount int
}
