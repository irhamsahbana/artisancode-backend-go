package restentity

type AttendanceLog struct {
	ID             string   `json:"id"`
	EmployeeID     string   `json:"employee_id"`
	EmployeeNo     string   `json:"employee_no"`
	EmployeeName   string   `json:"employee_name"`
	AttendanceDate string   `json:"attendance_date"`
	Type           string   `json:"type"`
	Source         string   `json:"source"`
	Status         string   `json:"status"`
	LoggedAt       string   `json:"logged_at"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	Address        *string  `json:"address"`
	DeviceID       *string  `json:"device_id"`
	DeviceName     *string  `json:"device_name"`
	Notes          *string  `json:"notes"`
	SelfieFileID   *string  `json:"selfie_file_id"`
	SelfieURL      *string  `json:"selfie_url"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      *string  `json:"updated_at"`
}

type GetAttendanceLogsReq struct {
	Q              string  `query:"q" validate:"omitempty,min=2"`
	EmployeeID     *string `query:"employee_id" validate:"omitempty,uuidv7"`
	Type           string  `query:"type" validate:"omitempty,oneof=check_in check_out"`
	Source         string  `query:"source" validate:"omitempty,oneof=mobile web"`
	Status         string  `query:"status" validate:"omitempty,oneof=recorded"`
	SelfieStatus   string  `query:"selfie_status" validate:"omitempty,oneof=with_photo without_photo"`
	OrgUnitID      *string `query:"org_unit_id" validate:"omitempty,uuidv7"`
	BranchID       *string `query:"branch_id" validate:"omitempty,uuidv7"`
	WorkLocationID *string `query:"work_location_id" validate:"omitempty,uuidv7"`
	ExceptionType  string  `query:"exception_type" validate:"omitempty,oneof=late_check_in missing_check_out missing_check_in"`
	AttendanceDay  *string `query:"attendance_date" validate:"omitempty,datetime=2006-01-02"`
	DateFrom       *string `query:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo         *string `query:"date_to" validate:"omitempty,datetime=2006-01-02"`
	Page           int     `query:"page" validate:"omitempty,min=1"`
	Limit          int     `query:"limit" validate:"omitempty,min=1"`
}

func (r *GetAttendanceLogsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 15
	}
}

type GetAttendanceLogsResp struct {
	Items      []AttendanceLog `json:"items"`
	Pagination PaginationResp  `json:"pagination"`
}

type GetAttendanceLogReq struct {
	ID string `params:"id" validate:"required"`
}

type GetAttendanceLogResp struct {
	AttendanceLog
}

type CheckAttendanceReq struct {
	LoggedAt     string   `json:"logged_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Address      *string  `json:"address"`
	DeviceID     *string  `json:"device_id" validate:"omitempty,max=255"`
	DeviceName   *string  `json:"device_name" validate:"omitempty,max=255"`
	Notes        *string  `json:"notes" validate:"omitempty,max=1000"`
	SelfieFileID string   `json:"selfie_file_id" validate:"required,uuidv7"`
}

type CheckAttendanceResp struct {
	ID string `json:"id"`
}

type GetAttendanceSummaryTodayResp struct {
	AttendanceDate string  `json:"attendance_date"`
	CheckedIn      bool    `json:"checked_in"`
	CheckedOut     bool    `json:"checked_out"`
	CheckInLogID   *string `json:"check_in_log_id"`
	CheckOutLogID  *string `json:"check_out_log_id"`
	LastLogType    *string `json:"last_log_type"`
	LastLoggedAt   *string `json:"last_logged_at"`
	CanCheckIn     bool    `json:"can_check_in"`
	CanCheckOut    bool    `json:"can_check_out"`
}

type GetAttendancePolicyResp struct {
	Timezone                string `json:"timezone"`
	AttendanceRadiusMeters  int    `json:"attendance_radius_meters"`
	AttendanceCheckInStart  string `json:"attendance_check_in_start"`
	AttendanceCheckInEnd    string `json:"attendance_check_in_end"`
	AttendanceCheckOutStart string `json:"attendance_check_out_start"`
	AttendanceCheckOutEnd   string `json:"attendance_check_out_end"`
}

type GetOwnerAttendanceDashboardReq struct {
	Date      *string `query:"date" validate:"omitempty,datetime=2006-01-02"`
	Timezone  string  `query:"timezone" validate:"required,timezone"`
	TrendDays int     `query:"trend_days" validate:"omitempty,oneof=7 14 30"`
}

func (r *GetOwnerAttendanceDashboardReq) SetDefault() {
	if r.Date == nil || *r.Date == "" {
		today := TodayDate()
		r.Date = &today
	}
	if r.TrendDays == 0 {
		r.TrendDays = 7
	}
}

type GetOwnerAttendanceDashboardResp struct {
	AttendanceDate  string                                   `json:"attendance_date"`
	Summary         OwnerAttendanceDashboardSummaryResp      `json:"summary"`
	TodayExceptions []OwnerAttendanceDashboardExceptionResp  `json:"today_exceptions"`
	DailyTrend      []OwnerAttendanceDashboardDailyTrendResp `json:"daily_trend"`
}

type OwnerAttendanceDashboardSummaryResp struct {
	ActiveEmployeeCount  int `json:"active_employee_count"`
	CheckedInCount       int `json:"checked_in_count"`
	CheckedOutCount      int `json:"checked_out_count"`
	PendingCheckInCount  int `json:"pending_check_in_count"`
	PendingCheckOutCount int `json:"pending_check_out_count"`
	LateCheckInCount     int `json:"late_check_in_count"`
}

type OwnerAttendanceDashboardExceptionResp struct {
	EmployeeID     string  `json:"employee_id"`
	EmployeeNo     string  `json:"employee_no"`
	EmployeeName   string  `json:"employee_name"`
	ShiftName      *string `json:"shift_name"`
	FirstCheckInAt *string `json:"first_check_in_at"`
	LastCheckOutAt *string `json:"last_check_out_at"`
	ExceptionType  string  `json:"exception_type"`
}

type OwnerAttendanceDashboardDailyTrendResp struct {
	AttendanceDate       string `json:"attendance_date"`
	CheckedInCount       int    `json:"checked_in_count"`
	CheckedOutCount      int    `json:"checked_out_count"`
	LateCheckInCount     int    `json:"late_check_in_count"`
	MissingCheckOutCount int    `json:"missing_check_out_count"`
}
