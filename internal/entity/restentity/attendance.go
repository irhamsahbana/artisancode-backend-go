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
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      *string  `json:"updated_at"`
}

type GetAttendanceLogsReq struct {
	Q             string  `query:"q" validate:"omitempty,min=2"`
	EmployeeID    *string `query:"employee_id" validate:"omitempty,uuidv7"`
	Type          string  `query:"type" validate:"omitempty,oneof=check_in check_out"`
	Source        string  `query:"source" validate:"omitempty,oneof=mobile web"`
	AttendanceDay *string `query:"attendance_date" validate:"omitempty,datetime=2006-01-02"`
	DateFrom      *string `query:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo        *string `query:"date_to" validate:"omitempty,datetime=2006-01-02"`
	Page          int     `query:"page" validate:"omitempty,min=1"`
	Limit         int     `query:"limit" validate:"omitempty,min=1"`
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
	LoggedAt   string   `json:"logged_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Address    *string  `json:"address"`
	DeviceID   *string  `json:"device_id" validate:"omitempty,max=255"`
	DeviceName *string  `json:"device_name" validate:"omitempty,max=255"`
	Notes      *string  `json:"notes" validate:"omitempty,max=1000"`
}

type CheckAttendanceResp struct {
	ID string `json:"id"`
}
