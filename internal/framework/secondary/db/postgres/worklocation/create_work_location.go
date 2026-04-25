package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *workLocationRepo) CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:CreateWorkLocation")
	defer span.End()

	query := `
		INSERT INTO work_locations (
			tenant_id, org_unit_id, name, address, latitude, longitude, radius_meters
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.OrgUnitID, data.Name, data.Address, data.Latitude, data.Longitude, data.RadiusMeters,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create work location")
		return nil, err
	}

	data.ID = id
	return &data, nil
}
