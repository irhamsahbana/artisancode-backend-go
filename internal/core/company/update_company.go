package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *companyCore) UpdateCompany(ctx context.Context, data coreentity.Company) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:company:update_company:UpdateCompany")
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

	return c.repo.UpdateCompany(ctx, data)
}
