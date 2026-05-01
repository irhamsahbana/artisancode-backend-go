package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *workLocationCore) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:worklocation:update_work_location:UpdateWorkLocation")
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
		log.Ctx(ctx).Warn().Any("name", data.Name).Msg(errmsg.MessageWorkLocationNameAlreadyExists)
		return errmsg.NewCustomErrors(409).SetMessage(errmsg.MessageWorkLocationNameAlreadyExists)
	}

	// Validate org_unit_id exists if provided
	if data.OrgUnitID != nil {
		_, err := c.orgUnitRepo.GetOrgUnit(ctx, coreentity.OrgUnit{
			TenantID: data.TenantID,
			ID:       *data.OrgUnitID,
		})
		if err != nil {
			log.Ctx(ctx).Warn().Str("org_unit_id", *data.OrgUnitID).Msg(errmsg.MessageOrgUnitNotFound)
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageOrganizationUnitNotFound)
		}
	}

	return c.repo.UpdateWorkLocation(ctx, data)
}
