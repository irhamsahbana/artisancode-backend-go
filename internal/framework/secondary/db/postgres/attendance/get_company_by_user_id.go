package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetCompanyByUserID(ctx context.Context, tenantID, userID string) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:attendance:get_company_by_user_id:GetCompanyByUserID")
	defer span.End()

	var data struct {
		ID       string `db:"id"`
		TenantID string `db:"tenant_id"`
		Name     string `db:"name"`
		Code     string `db:"code"`
		Config   []byte `db:"config"`
	}

	query := `
		WITH RECURSIVE org_unit_ancestors AS (
			SELECT ou.id, ou.tenant_id, ou.name, ou.code, ou.parent_id, ou.category, ou.config
			FROM employees e
			INNER JOIN org_units ou ON ou.id = e.org_unit_id AND ou.deleted_at IS NULL
			WHERE e.tenant_id = ? AND e.user_id = ? AND e.deleted_at IS NULL

			UNION ALL

			SELECT parent.id, parent.tenant_id, parent.name, parent.code, parent.parent_id, parent.category, parent.config
			FROM org_units parent
			INNER JOIN org_unit_ancestors child ON child.parent_id = parent.id
			WHERE parent.deleted_at IS NULL
		)
		SELECT id, tenant_id, name, code, config
		FROM org_unit_ancestors
		WHERE category = 'company'
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), tenantID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"user_id":   userID,
		}).Msg("Failed to resolve company by user id")
		return nil, err
	}

	var config coreentity.CompanyConfig
	if len(data.Config) > 0 {
		err = json.Unmarshal(data.Config, &config)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data.ID).Msg("Failed to unmarshal company config")
			return nil, err
		}
	}

	return &coreentity.Company{
		ID:       data.ID,
		TenantID: data.TenantID,
		Name:     data.Name,
		Code:     data.Code,
		Config:   config,
	}, nil
}
