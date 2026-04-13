package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type workLocationRepo struct {
	db *sqlx.DB
}

var _ portsRepo.WorkLocationRepository = &workLocationRepo{}

type WorkLocationRepositoryConfig struct {
	DB *sqlx.DB
}

func NewWorkLocationRepository(cfg WorkLocationRepositoryConfig) portsRepo.WorkLocationRepository {
	return &workLocationRepo{db: cfg.DB}
}

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
			wl.name, wl.address, wl.timezone, wl.latitude, wl.longitude, wl.radius_meters
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
			Timezone:     d.Timezone,
			Latitude:     d.Latitude,
			Longitude:    d.Longitude,
			RadiusMeters: d.RadiusMeters,
		})
	}

	return items, total, nil
}

func (r *workLocationRepo) GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:GetWorkLocation")
	defer span.End()

	var data = new(repoentity.WorkLocation)

	query := `
		SELECT wl.id, wl.tenant_id, wl.org_unit_id, ou.name AS org_unit_name,
			wl.name, wl.address, wl.timezone, wl.latitude, wl.longitude, wl.radius_meters
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
		Timezone:     data.Timezone,
		Latitude:     data.Latitude,
		Longitude:    data.Longitude,
		RadiusMeters: data.RadiusMeters,
	}
	return result, nil
}

func (r *workLocationRepo) CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:CreateWorkLocation")
	defer span.End()

	query := `
		INSERT INTO work_locations (
			tenant_id, org_unit_id, name, address, timezone, latitude, longitude, radius_meters
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.OrgUnitID, data.Name, data.Address, data.Timezone, data.Latitude, data.Longitude, data.RadiusMeters,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create work location")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *workLocationRepo) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:UpdateWorkLocation")
	defer span.End()

	query := `
		UPDATE work_locations
		SET org_unit_id = ?, name = ?, address = ?, timezone = ?, latitude = ?, longitude = ?, radius_meters = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.OrgUnitID, data.Name, data.Address, data.Timezone, data.Latitude, data.Longitude, data.RadiusMeters, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update work location")
		return err
	}
	return nil
}

func (r *workLocationRepo) DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:DeleteWorkLocation")
	defer span.End()

	query := `
		UPDATE work_locations
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete work location")
		return err
	}
	return nil
}

func (r *workLocationRepo) ExistsWorkLocationByName(ctx context.Context, tenantID string, name string, excludeID *string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:ExistsWorkLocationByName")
	defer span.End()

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM work_locations
			WHERE tenant_id = ? AND name = ? AND deleted_at IS NULL
		`
	args := []any{tenantID, name}

	if excludeID != nil {
		query += ` AND id != ?`
		args = append(args, *excludeID)
	}

	query += `)`

	err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("name", name).Msg("Failed to check work location name existence")
		return false, err
	}

	return exists, nil
}
