package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/repository"
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
	ctx, span := tracing.StartSpan(ctx, "repo.GetWorkLocations")
	defer span.End()

	type dao struct {
		TotalData int `db:"total_data"`
		repoentity.WorkLocation
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 4)
		items = make([]coreentity.WorkLocation, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, address, timezone, latitude, longitude, radius_meters
		FROM work_locations
		WHERE deleted_at IS NULL AND tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
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
	ctx, span := tracing.StartSpan(ctx, "repo.GetWorkLocation")
	defer span.End()

	var data = new(repoentity.WorkLocation)

	query := `
		SELECT id, tenant_id, name, address, timezone, latitude, longitude, radius_meters
		FROM work_locations
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
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
	ctx, span := tracing.StartSpan(ctx, "repo.CreateWorkLocation")
	defer span.End()

	query := `
		INSERT INTO work_locations (
			tenant_id, name, address, timezone, latitude, longitude, radius_meters
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.Name, data.Address, data.Timezone, data.Latitude, data.Longitude, data.RadiusMeters,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create work location")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *workLocationRepo) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "repo.UpdateWorkLocation")
	defer span.End()

	query := `
		UPDATE work_locations
		SET name = ?, address = ?, timezone = ?, latitude = ?, longitude = ?, radius_meters = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Name, data.Address, data.Timezone, data.Latitude, data.Longitude, data.RadiusMeters, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update work location")
		return err
	}
	return nil
}

func (r *workLocationRepo) DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo.DeleteWorkLocation")
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
