package core

import (
	"context"
	"slices"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

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
	if cfg.LeaveAllowanceAnnual < 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Leave allowance annual must be >= 0")
	}
	if cfg.OvertimeRateMultiplier < 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Overtime rate multiplier must be >= 0")
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
	if cfg.PreferredLanguage == "" {
		return errmsg.NewCustomErrors(400).SetMessage("Preferred language is required")
	}
	if len(cfg.SupportedLanguages) == 0 {
		return errmsg.NewCustomErrors(400).SetMessage("Supported languages is required")
	}

	allowedLanguages := []string{"id", "en"}
	for _, language := range cfg.SupportedLanguages {
		if !slices.Contains(allowedLanguages, language) {
			return errmsg.NewCustomErrors(400).SetMessage("Supported languages contains unsupported language")
		}
	}
	if !slices.Contains(cfg.SupportedLanguages, cfg.PreferredLanguage) {
		return errmsg.NewCustomErrors(400).SetMessage("Preferred language must exist in supported languages")
	}

	return nil
}
