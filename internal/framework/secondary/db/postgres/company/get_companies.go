package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

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
