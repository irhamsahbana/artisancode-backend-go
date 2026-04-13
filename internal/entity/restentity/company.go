package restentity

import "codebase-app/pkg/types"

type Company struct {
	ID        string        `json:"id"`
	Code      string        `json:"code"`
	Name      string        `json:"name"`
	Config    CompanyConfig `json:"config"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

type CompanyConfig struct {
	Logo                    *string  `json:"logo"`
	AttendanceRadiusMeters  int      `json:"attendance_radius_meters"`
	AttendanceCheckInStart  string   `json:"attendance_check_in_start"`
	AttendanceCheckInEnd    string   `json:"attendance_check_in_end"`
	AttendanceCheckOutStart string   `json:"attendance_check_out_start"`
	AttendanceCheckOutEnd   string   `json:"attendance_check_out_end"`
	LeaveAllowanceAnnual    int      `json:"leave_allowance_annual"`
	OvertimeRateMultiplier  float64  `json:"overtime_rate_multiplier"`
	PreferredLanguage       string   `json:"preferred_language"`
	SupportedLanguages      []string `json:"supported_languages"`
	Timezone                string   `json:"timezone"`
	DateFormat              string   `json:"date_format"`
	TimeFormat              string   `json:"time_format"`
}

type GetCompaniesReq struct {
	Q string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetCompaniesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetCompaniesResp struct {
	Items []Company  `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type GetCompanyReq struct {
	ID string `params:"id" validate:"required"`
}

type GetCompanyResp struct {
	Company
}

type CreateCompanyReq struct {
	Code   string        `json:"code" validate:"required,min=2,max=64"`
	Name   string        `json:"name" validate:"required,min=2"`
	Config CompanyConfig `json:"config" validate:"required"`
}

type CreateCompanyResp struct {
	ID string `json:"id"`
}

type UpdateCompanyReq struct {
	ID     string        `params:"id" validate:"required"`
	Code   string        `json:"code" validate:"required,min=2,max=64"`
	Name   string        `json:"name" validate:"required,min=2"`
	Config CompanyConfig `json:"config" validate:"required"`
}

type DeleteCompanyReq struct {
	ID string `params:"id" validate:"required"`
}
