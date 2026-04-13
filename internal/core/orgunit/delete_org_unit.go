package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
)

func (c *orgUnitCore) DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:delete_org_unit:DeleteOrgUnit")
	defer span.End()

	// Check if org unit exists before deleting
	_, err := c.repo.GetOrgUnit(ctx, coreentity.OrgUnit{
		TenantID: filter.TenantID,
		ID:       filter.ID,
	})
	if err != nil {
		return err
	}

	return c.repo.DeleteOrgUnit(ctx, filter)
}
