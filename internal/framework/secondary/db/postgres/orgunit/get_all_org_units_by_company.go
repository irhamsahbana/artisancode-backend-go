package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

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
