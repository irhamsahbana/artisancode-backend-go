package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *orgUnitRepo) GetOrgUnits(
	ctx context.Context,
	filter coreentity.OrgUnitListFilter,
) ([]coreentity.OrgUnit, int, error) {
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
