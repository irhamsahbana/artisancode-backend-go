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

	if c.billingCore != nil {
		count, err := c.repo.CountCompaniesByTenant(ctx, data.TenantID)
		if err != nil {
			return nil, err
		}
		if err := c.billingCore.CheckUsageLimit(ctx, coreentity.ResourceTypeBranches, count); err != nil {
			return nil, err
		}
	}

	// Validate code uniqueness
	exists, err := c.repo.ExistsCompanyByCode(ctx, data.TenantID, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"code": data.Code,
		}).Msg(errmsg.MessageCompanyCodeAlreadyExists)
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCompanyCodeAlreadyExists)
	}

	// Validate config
	if err := validateCompanyConfig(data.Config); err != nil {
		return nil, err
	}

	return c.repo.CreateCompany(ctx, data)
}

func validateCompanyConfig(cfg coreentity.CompanyConfig) error {
	if cfg.Timezone == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageTimezoneIsRequired)
	}
	if cfg.DateFormat == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageDateFormatIsRequired)
	}
	if cfg.TimeFormat == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageTimeFormatIsRequired)
	}
	if cfg.PreferredLanguage == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePreferredLanguageIsRequired)
	}
	if len(cfg.SupportedLanguages) == 0 {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageSupportedLanguagesIsRequired)
	}

	allowedLanguages := []string{"id", "en"}
	for _, language := range cfg.SupportedLanguages {
		if !slices.Contains(allowedLanguages, language) {
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageSupportedLanguagesContainsUnsupportedLanguage)
		}
	}
	if !slices.Contains(cfg.SupportedLanguages, cfg.PreferredLanguage) {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePreferredLanguageMustExistInSupportedLanguages)
	}

	return nil
}
