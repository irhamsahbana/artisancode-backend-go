package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
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
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:GetOrgUnits")
	defer span.End()

	type dao struct {
		TotalData int     `db:"total_data"`
		ID        string  `db:"id"`
		TenantID  string  `db:"tenant_id"`
		Code      string  `db:"code"`
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
			id, tenant_id, code, name, parent_id, category
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
			Code:     d.Code,
			Name:     d.Name,
			ParentID: d.ParentID,
			Category: d.Category,
		})
	}

	return items, total, nil
}

func (r *orgUnitRepo) GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:GetOrgUnit")
	defer span.End()

	var data struct {
		ID       string  `db:"id"`
		TenantID string  `db:"tenant_id"`
		Code     string  `db:"code"`
		Name     string  `db:"name"`
		ParentID *string `db:"parent_id"`
		Category string  `db:"category"`
	}

	query := `
		SELECT id, tenant_id, code, name, parent_id, category
		FROM org_units
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Org unit not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage("Org unit not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get org unit")
		return nil, err
	}

	result := &coreentity.OrgUnit{
		ID:       data.ID,
		TenantID: data.TenantID,
		Code:     data.Code,
		Name:     data.Name,
		ParentID: data.ParentID,
		Category: data.Category,
	}
	return result, nil
}

func (r *orgUnitRepo) CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:CreateOrgUnit")
	defer span.End()

	query := `
		INSERT INTO org_units (tenant_id, code, name, parent_id, category)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`

	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query), data.TenantID, data.Code, data.Name, data.ParentID, data.Category)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create org unit")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *orgUnitRepo) UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:UpdateOrgUnit")
	defer span.End()

	query := `
		UPDATE org_units
		SET code = ?, name = ?, parent_id = ?, category = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.Code, data.Name, data.ParentID, data.Category, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update org unit")
		return err
	}
	return nil
}

func (r *orgUnitRepo) GetAllOrgUnitsByCompany(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:GetAllOrgUnitsByCompany")
	defer span.End()

	payload := map[string]any{
		"tenant_id":  tenantID,
		"company_id": companyID,
	}

	type dao struct {
		ID       string  `db:"id"`
		TenantID string  `db:"tenant_id"`
		Code     string  `db:"code"`
		Name     string  `db:"name"`
		ParentID *string `db:"parent_id"`
		Category string  `db:"category"`
	}

	var (
		data  = make([]dao, 0)
		items = make([]coreentity.OrgUnit, 0)
	)

	// Recursive CTE to get all children of company org unit
	query := `
		WITH RECURSIVE org_unit_hierarchy AS (
			SELECT id, tenant_id, code, name, parent_id, category
			FROM org_units
			WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL

			UNION ALL

			SELECT ou.id, ou.tenant_id, ou.code, ou.name, ou.parent_id, ou.category
			FROM org_units ou
			INNER JOIN org_unit_hierarchy oh ON ou.parent_id = oh.id
			WHERE ou.deleted_at IS NULL
		)
		SELECT id, tenant_id, code, name, parent_id, category
		FROM org_unit_hierarchy
		ORDER BY name ASC
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), companyID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to query all org units for company")
		return nil, err
	}

	for _, d := range data {
		items = append(items, coreentity.OrgUnit{
			ID:       d.ID,
			TenantID: d.TenantID,
			Code:     d.Code,
			Name:     d.Name,
			ParentID: d.ParentID,
			Category: d.Category,
		})
	}

	return items, nil
}

func (r *orgUnitRepo) ExistsOrgUnitByCode(ctx context.Context, tenantID string, code string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:ExistsOrgUnitByCode")
	defer span.End()

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM org_units
			WHERE tenant_id = ? AND code = ? AND deleted_at IS NULL
		)
	`

	err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), tenantID, code)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID, "code": code}).
			Msg("Failed to check org unit code existence")
		return false, err
	}

	return exists, nil
}

func (r *orgUnitRepo) DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:DeleteOrgUnit")
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
