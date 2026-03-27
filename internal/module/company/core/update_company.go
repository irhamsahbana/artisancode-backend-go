package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *companyCore) UpdateCompany(ctx context.Context, data coreentity.Company) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateCompany")
	defer span.End()

	// Verify company exists
	existing, err := c.repo.GetCompany(ctx, coreentity.Company{
		TenantID: data.TenantID,
		ID:       data.ID,
	})
	if err != nil {
		return err
	}
	if existing == nil {
		return errmsg.NewCustomErrors(404).SetMessage("Company not found")
	}

	// Validate code uniqueness (exclude current company)
	exists, err := c.repo.ExistsCompanyByCode(ctx, data.TenantID, data.Code, data.ID)
	if err != nil {
		return err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"code": data.Code,
		}).Msg("Company code already exists")
		return errmsg.NewCustomErrors(400).SetMessage("Company code already exists")
	}

	// Validate config
	if err := validateCompanyConfig(data.Config); err != nil {
		return err
	}

	return c.repo.UpdateCompany(ctx, data)
}