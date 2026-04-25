package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) GetPermissions(ctx context.Context, filter coreentity.PermissionListFilter) ([]coreentity.Permission, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetPermissions")
	defer span.End()

	type row struct {
		TotalData int `db:"total_data"`
		coreentity.Permission
	}

	rows := make([]row, 0)
	args := make([]any, 0, 4)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			name,
			COALESCE(description, '') AS description
		FROM permissions
		WHERE deleted_at IS NULL
	`
	if filter.TenantID != "" {
		query += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		return nil, 0, err
	}

	items := make([]coreentity.Permission, 0, len(rows))
	total := 0
	for _, item := range rows {
		total = item.TotalData
		items = append(items, item.Permission)
	}
	return items, total, nil
}
