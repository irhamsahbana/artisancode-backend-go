package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/repository"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type orgUnitRepo struct {
	db *sqlx.DB
}

var _ portsRepo.OrgUnitRepository = &orgUnitRepo{}

type OrgUnitRepositoryConfig struct {
	DB *sqlx.DB
}

func NewOrgUnitRepository(cfg OrgUnitRepositoryConfig) portsRepo.OrgUnitRepository {
	return &orgUnitRepo{db: cfg.DB}
}

func (r *orgUnitRepo) GetOrgUnits(ctx context.Context, filter coreentity.OrgUnitListFilter) ([]coreentity.OrgUnit, int, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetOrgUnits")
	defer span.End()

	type dao struct {
		TotalData int     `db:"total_data"`
		ID        string  `db:"id"`
		TenantID  string  `db:"tenant_id"`
		Name      string  `db:"name"`
		ParentID  *string `db:"parent_id"`
		Category  string  `db:"category"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 6)
		items = make([]coreentity.OrgUnit, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, parent_id, category
		FROM org_units
		WHERE deleted_at IS NULL AND tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.Category != "" {
		query += ` AND category = ?`
		args = append(args, filter.Category)
	}
	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query org units")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.OrgUnit{
			ID:       d.ID,
			TenantID: d.TenantID,
			Name:     d.Name,
			ParentID: d.ParentID,
			Category: d.Category,
		})
	}

	return items, total, nil
}

func (r *orgUnitRepo) GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetOrgUnit")
	defer span.End()

	var data struct {
		ID       string  `db:"id"`
		TenantID string  `db:"tenant_id"`
		Name     string  `db:"name"`
		ParentID *string `db:"parent_id"`
		Category string  `db:"category"`
	}

	query := `
		SELECT id, tenant_id, name, parent_id, category
		FROM org_units
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Org unit tidak ditemukan")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get org unit")
		return nil, err
	}

	result := &coreentity.OrgUnit{
		ID:       data.ID,
		TenantID: data.TenantID,
		Name:     data.Name,
		ParentID: data.ParentID,
		Category: data.Category,
	}
	return result, nil
}

func (r *orgUnitRepo) CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateOrgUnit")
	defer span.End()

	query := `
		INSERT INTO org_units (tenant_id, name, parent_id, category)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`

	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query), data.TenantID, data.Name, data.ParentID, data.Category)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create org unit")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *orgUnitRepo) UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error {
	ctx, span := tracing.StartSpan(ctx, "repo.UpdateOrgUnit")
	defer span.End()

	query := `
		UPDATE org_units
		SET name = ?, parent_id = ?, category = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.Name, data.ParentID, data.Category, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update org unit")
		return err
	}
	return nil
}

func (r *orgUnitRepo) DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo.DeleteOrgUnit")
	defer span.End()

	query := `
		UPDATE org_units
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete org unit")
		return err
	}
	return nil
}
