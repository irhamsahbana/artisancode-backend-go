package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *workLocationCore) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateWorkLocation")
	defer span.End()

	// Validate work location exists
	_, err := c.repo.GetWorkLocation(ctx, coreentity.WorkLocation{
		TenantID: data.TenantID,
		ID:       data.ID,
	})
	if err != nil {
		return err
	}

	// Validate name uniqueness (exclude self)
	exists, err := c.repo.ExistsWorkLocationByName(ctx, data.TenantID, data.Name, &data.ID)
	if err != nil {
		return err
	}
	if exists {
		log.Ctx(ctx).Warn().Any("name", data.Name).Msg("Work location name already exists")
		return errmsg.NewCustomErrors(409).SetMessage("Work location name already exists")
	}

	// Validate org_unit_id exists if provided
	if data.OrgUnitID != nil {
		_, err := c.orgUnitRepo.GetOrgUnit(ctx, coreentity.OrgUnit{
			TenantID: data.TenantID,
			ID:       *data.OrgUnitID,
		})
		if err != nil {
			log.Ctx(ctx).Warn().Str("org_unit_id", *data.OrgUnitID).Msg("Org unit not found")
			return errmsg.NewCustomErrors(400).SetMessage("Organization unit not found")
		}
	}

	return c.repo.UpdateWorkLocation(ctx, data)
}