package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *workLocationRepo) GetWorkLocations(ctx context.Context, filter coreentity.WorkLocationListFilter) ([]coreentity.WorkLocation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:GetWorkLocations")
	defer span.End()

	type dao struct {
		TotalData int `db:"total_data"`
		repoentity.WorkLocation
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 5)
		items = make([]coreentity.WorkLocation, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			wl.id, wl.tenant_id, wl.org_unit_id, ou.name AS org_unit_name,
			wl.name, wl.address, wl.latitude, wl.longitude, wl.radius_meters
		FROM work_locations wl
		LEFT JOIN org_units ou ON wl.org_unit_id = ou.id
		WHERE wl.deleted_at IS NULL AND wl.tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.OrgUnitID != nil {
		query += ` AND wl.org_unit_id = ?`
		args = append(args, *filter.OrgUnitID)
	}

	if filter.Q != "" {
		query += ` AND wl.name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY wl.name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query work locations")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.WorkLocation{
			ID:           d.ID,
			TenantID:     d.TenantID,
			OrgUnitID:    d.OrgUnitID,
			OrgUnitName:  d.OrgUnitName,
			Name:         d.Name,
			Address:      d.Address,
			Latitude:     d.Latitude,
			Longitude:    d.Longitude,
			RadiusMeters: d.RadiusMeters,
		})
	}

	return items, total, nil
}
