package restentity

type CreateExportJobReq struct {
	ResourceType   string  `json:"resource_type" validate:"required,oneof=attendance_logs"`
	Format         string  `json:"format" validate:"required,oneof=xlsx"`
	Q              string  `json:"q" validate:"omitempty,min=2"`
	EmployeeID     *string `json:"employee_id" validate:"omitempty,uuidv7"`
	Type           string  `json:"type" validate:"omitempty,oneof=check_in check_out"`
	Source         string  `json:"source" validate:"omitempty,oneof=mobile web"`
	Status         string  `json:"status" validate:"omitempty,oneof=recorded"`
	SelfieStatus   string  `json:"selfie_status" validate:"omitempty,oneof=with_photo without_photo"`
	OrgUnitID      *string `json:"org_unit_id" validate:"omitempty,uuidv7"`
	BranchID       *string `json:"branch_id" validate:"omitempty,uuidv7"`
	WorkLocationID *string `json:"work_location_id" validate:"omitempty,uuidv7"`
	ExceptionType  string  `json:"exception_type" validate:"omitempty,oneof=late_check_in missing_check_out missing_check_in"`
	AttendanceDay  *string `json:"attendance_date" validate:"omitempty,datetime=2006-01-02"`
	DateFrom       *string `json:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo         *string `json:"date_to" validate:"omitempty,datetime=2006-01-02"`
}

type CreateExportJobResp struct {
	ID string `json:"id"`
}

type GetExportJobsReq struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1"`
}

func (r *GetExportJobsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 10
	}
}

type ExportJob struct {
	ID              string  `json:"id"`
	ResourceType    string  `json:"resource_type"`
	ResourceLabel   string  `json:"resource_label"`
	Format          string  `json:"format"`
	Status          string  `json:"status"`
	RequestedByName string  `json:"requested_by_name"`
	ErrorMessage    *string `json:"error_message"`
	StartedAt       *string `json:"started_at"`
	CompletedAt     *string `json:"completed_at"`
	ExpiresAt       *string `json:"expires_at"`
	CreatedAt       string  `json:"created_at"`
	DownloadURL     *string `json:"download_url"`
}

type GetExportJobsResp struct {
	Items      []ExportJob    `json:"items"`
	Pagination PaginationResp `json:"pagination"`
}

type GetExportJobReq struct {
	ID string `params:"id" validate:"required,uuidv7"`
}

type GetExportJobResp struct {
	ExportJob
}
