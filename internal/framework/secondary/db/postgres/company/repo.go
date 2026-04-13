package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type companyRepo struct {
	db *sqlx.DB
}

var _ portsRepo.CompanyRepository = &companyRepo{}

type CompanyRepositoryConfig struct {
	DB *sqlx.DB
}

func NewCompanyRepository(cfg CompanyRepositoryConfig) portsRepo.CompanyRepository {
	return &companyRepo{db: cfg.DB}
}

func (r *companyRepo) GetCompanies(ctx context.Context, filter coreentity.CompanyListFilter) ([]coreentity.Company, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:GetCompanies")
	defer span.End()

	type dao struct {
		TotalData int     `db:"total_data"`
		ID        string  `db:"id"`
		TenantID  string  `db:"tenant_id"`
		Name      string  `db:"name"`
		Code      string  `db:"code"`
		Config    []byte  `db:"config"`
		CreatedAt string  `db:"created_at"`
		UpdatedAt *string `db:"updated_at"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 4)
		items = make([]coreentity.Company, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, code, config, created_at, updated_at
		FROM org_units
		WHERE deleted_at IS NULL AND tenant_id = ? AND category = 'company'
	`

	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query companies")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		var config coreentity.CompanyConfig
		if len(d.Config) > 0 {
			err = json.Unmarshal(d.Config, &config)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, d.ID).Msg("Failed to unmarshal config")
				return nil, 0, err
			}
		}
		updatedAt := ""
		if d.UpdatedAt != nil {
			updatedAt = *d.UpdatedAt
		}

		items = append(items, coreentity.Company{
			ID:        d.ID,
			TenantID:  d.TenantID,
			Name:      d.Name,
			Code:      d.Code,
			Config:    config,
			CreatedAt: d.CreatedAt,
			UpdatedAt: updatedAt,
		})
	}

	return items, total, nil
}

func (r *companyRepo) GetCompany(ctx context.Context, filter coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:GetCompany")
	defer span.End()

	type dao struct {
		ID        string  `db:"id"`
		TenantID  string  `db:"tenant_id"`
		Name      string  `db:"name"`
		Code      string  `db:"code"`
		Config    []byte  `db:"config"`
		CreatedAt string  `db:"created_at"`
		UpdatedAt *string `db:"updated_at"`
	}

	query := `
		SELECT id, tenant_id, name, code, config, created_at, updated_at
		FROM org_units
		WHERE id = ? AND tenant_id = ? AND category = 'company' AND deleted_at IS NULL
	`

	var data dao
	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Company not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get company")
		return nil, err
	}

	var config coreentity.CompanyConfig
	if len(data.Config) > 0 {
		err = json.Unmarshal(data.Config, &config)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data.ID).Msg("Failed to unmarshal config")
			return nil, err
		}
	}

	updatedAt := ""
	if data.UpdatedAt != nil {
		updatedAt = *data.UpdatedAt
	}

	result := &coreentity.Company{
		ID:        data.ID,
		TenantID:  data.TenantID,
		Name:      data.Name,
		Code:      data.Code,
		Config:    config,
		CreatedAt: data.CreatedAt,
		UpdatedAt: updatedAt,
	}
	return result, nil
}

func (r *companyRepo) CreateCompany(ctx context.Context, data coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:CreateCompany")
	defer span.End()

	configJSON, err := json.Marshal(data.Config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to marshal config")
		return nil, err
	}

	query := `
		INSERT INTO org_units (tenant_id, category, name, code, config)
		VALUES (?, 'company', ?, ?, ?)
		RETURNING id
	`

	var id string
	err = r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.Name, data.Code, configJSON,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create company")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *companyRepo) UpdateCompany(ctx context.Context, data coreentity.Company) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:UpdateCompany")
	defer span.End()

	configJSON, err := json.Marshal(data.Config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to marshal config")
		return err
	}

	query := `
		UPDATE org_units
		SET name = ?, code = ?, config = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND category = 'company' AND deleted_at IS NULL
	`

	_, err = r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Name, data.Code, configJSON, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update company")
		return err
	}
	return nil
}

func (r *companyRepo) ExistsCompanyByCode(ctx context.Context, tenantID string, code string, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:ExistsCompanyByCode")
	defer span.End()

	query := `
		SELECT COUNT(*) FROM org_units
		WHERE deleted_at IS NULL AND tenant_id = ? AND category = 'company' AND code = ?
	`
	args := []any{tenantID, code}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID, "code": code}).Msg("Failed to check company code existence")
		return false, err
	}

	return count > 0, nil
}

func (r *companyRepo) DeleteCompany(ctx context.Context, filter coreentity.CompanyDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:DeleteCompany")
	defer span.End()

	query := `
		UPDATE org_units
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND category = 'company' AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete company")
		return err
	}
	return nil
}
