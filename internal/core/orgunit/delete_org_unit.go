package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"context"
)

func (c *orgUnitCore) DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:delete_org_unit:DeleteOrgUnit")
	defer span.End()

	// Check if org unit exists before deleting
	existing, err := c.repo.GetOrgUnit(ctx, coreentity.OrgUnit{
		TenantID: filter.TenantID,
		ID:       filter.ID,
	})
	if err != nil {
		return err
	}
	if existing.Category == "company" {
		return errmsg.NewCustomErrors(403).SetMessage("Company cannot be deleted")
	}

	return c.repo.DeleteOrgUnit(ctx, filter)
}
