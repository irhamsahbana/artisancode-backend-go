package coreentity

import (
	"codebase-app/internal/entity/common"
)

type Company struct {
	UserCtx common.UserContext

	ID        string
	TenantID  string
	Name      string
	Code      string
	Config    CompanyConfig
	CreatedAt string
	UpdatedAt string
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

func DefaultCompanyConfig() CompanyConfig {
	return CompanyConfig{
		Logo:                    nil,
		AttendanceRadiusMeters:  50,
		AttendanceCheckInStart:  "08:00",
		AttendanceCheckInEnd:    "09:00",
		AttendanceCheckOutStart: "17:00",
		AttendanceCheckOutEnd:   "23:59",
		LeaveAllowanceAnnual:    12,
		OvertimeRateMultiplier:  1.5,
		PreferredLanguage:       "id",
		SupportedLanguages:      []string{"id", "en"},
		Timezone:                "Asia/Jakarta",
		DateFormat:              "YYYY-MM-DD",
		TimeFormat:              "HH:mm:ss",
	}
}

type CompanyListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type CompanyDeleteFilter struct {
	TenantID string
	ID       string
}
