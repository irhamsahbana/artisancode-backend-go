package restentity

type MeResp struct {
	UserID      string   `json:"user_id"`
	UserName    string   `json:"user_name"`
	TenantID    string   `json:"tenant_id"`
	TenantName  string   `json:"tenant_name"`
	Roles       []string `json:"roles"`
	CompanyID   *string  `json:"company_id"`
	CompanyName *string  `json:"company_name"`
}

type GetMyEmployeeResp struct {
	ID            string  `json:"id"`
	EmployeeNo    string  `json:"employee_no"`
	FullName      string  `json:"full_name"`
	Email         string  `json:"email"`
	UserID        *string `json:"user_id"`
	OrgUnitID     *string `json:"org_unit_id"`
	JobPositionID *string `json:"job_position_id"`
	LocationID    *string `json:"location_id"`
	ShiftID       *string `json:"shift_id"`
	Status        string  `json:"status"`
	JoinDate      *string `json:"join_date"`
}

type GetMyShiftTodayResp struct {
	ShiftID        string `json:"shift_id"`
	ShiftName      string `json:"shift_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Timezone       string `json:"timezone"`
	AttendanceDate string `json:"attendance_date"`
}
