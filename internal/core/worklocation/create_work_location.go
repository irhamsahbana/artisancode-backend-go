package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *workLocationCore) CreateWorkLocation(
	ctx context.Context,
	data coreentity.WorkLocation,
) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:worklocation:create_work_location:CreateWorkLocation")
	defer span.End()

	// Validate name uniqueness
	exists, err := c.repo.ExistsWorkLocationByName(ctx, data.TenantID, data.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Ctx(ctx).Warn().Any("name", data.Name).Msg(errmsg.MessageWorkLocationNameAlreadyExists)
		return nil, errmsg.NewCustomErrors(409).SetMessage(errmsg.MessageWorkLocationNameAlreadyExists)
	}

	// Validate org_unit_id exists if provided
	if data.OrgUnitID != nil {
		_, err := c.orgUnitRepo.GetOrgUnit(ctx, coreentity.OrgUnit{
			TenantID: data.TenantID,
			ID:       *data.OrgUnitID,
		})
		if err != nil {
			log.Ctx(ctx).Warn().Str("org_unit_id", *data.OrgUnitID).Msg(errmsg.MessageOrgUnitNotFound)
			return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageOrganizationUnitNotFound)
		}
	}

	return c.repo.CreateWorkLocation(ctx, data)
}
