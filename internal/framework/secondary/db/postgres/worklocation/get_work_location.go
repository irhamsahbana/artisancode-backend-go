package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *workLocationRepo) GetWorkLocation(
	ctx context.Context,
	filter coreentity.WorkLocation,
) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:GetWorkLocation")
	defer span.End()

	var data = new(repoentity.WorkLocation)

	query := `
		SELECT wl.id, wl.tenant_id, wl.org_unit_id, ou.name AS org_unit_name,
			wl.name, wl.address, wl.latitude, wl.longitude, wl.radius_meters
		FROM work_locations wl
		LEFT JOIN org_units ou ON wl.org_unit_id = ou.id
		WHERE wl.id = ? AND wl.tenant_id = ? AND wl.deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Work location not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get work location")
		return nil, err
	}

	result := &coreentity.WorkLocation{
		ID:           data.ID,
		TenantID:     data.TenantID,
		OrgUnitID:    data.OrgUnitID,
		OrgUnitName:  data.OrgUnitName,
		Name:         data.Name,
		Address:      data.Address,
		Latitude:     data.Latitude,
		Longitude:    data.Longitude,
		RadiusMeters: data.RadiusMeters,
	}
	return result, nil
}
