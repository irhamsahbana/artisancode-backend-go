package core

import (
	"context"
	"regexp"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

var timeFormatRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)

func (c *companyCore) CreateCompany(ctx context.Context, data coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:company:create_company:CreateCompany")
	defer span.End()

	// Validate code uniqueness
	exists, err := c.repo.ExistsCompanyByCode(ctx, data.TenantID, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"code": data.Code,
		}).Msg("Company code already exists")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Company code already exists")
	}

	// Validate config
	if err := validateCompanyConfig(data.Config); err != nil {
		return nil, err
	}

	return c.repo.CreateCompany(ctx, data)
}

func validateCompanyConfig(cfg coreentity.CompanyConfig) error {
	if cfg.AttendanceRadiusMeters < 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Attendance radius must be >= 0")
	}
	if cfg.LeaveAllowanceAnnual < 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Leave allowance annual must be >= 0")
	}
	if cfg.OvertimeRateMultiplier < 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Overtime rate multiplier must be >= 0")
	}

	timeFields := map[string]string{
		"attendance_check_in_start":  cfg.AttendanceCheckInStart,
		"attendance_check_in_end":    cfg.AttendanceCheckInEnd,
		"attendance_check_out_start": cfg.AttendanceCheckOutStart,
		"attendance_check_out_end":   cfg.AttendanceCheckOutEnd,
	}
	for fieldName, value := range timeFields {
		if !timeFormatRegex.MatchString(value) {
			return errmsg.NewCustomErrors(400).SetMessage(fieldName + " format must be HH:mm")
		}
	}

	if cfg.Timezone == "" {
		return errmsg.NewCustomErrors(400).SetMessage("Timezone is required")
	}
	if cfg.DateFormat == "" {
		return errmsg.NewCustomErrors(400).SetMessage("Date format is required")
	}
	if cfg.TimeFormat == "" {
		return errmsg.NewCustomErrors(400).SetMessage("Time format is required")
	}

	return nil
}
