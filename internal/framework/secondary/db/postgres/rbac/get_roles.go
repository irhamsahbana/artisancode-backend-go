package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) GetRoles(ctx context.Context, filter coreentity.RoleListFilter) ([]coreentity.Role, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetRoles")
	defer span.End()

	type row struct {
		TotalData int     `db:"total_data"`
		ID        string  `db:"id"`
		TenantID  *string `db:"tenant_id"`
		Name      string  `db:"name"`
	}

	rows := make([]row, 0)
	args := make([]any, 0, 4)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			tenant_id,
			name
		FROM roles
		WHERE deleted_at IS NULL AND tenant_id = ?
	`
	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		return nil, 0, err
	}

	items := make([]coreentity.Role, 0, len(rows))
	total := 0
	for _, item := range rows {
		total = item.TotalData
		items = append(items, coreentity.Role{
			ID:       item.ID,
			TenantID: item.TenantID,
			Name:     item.Name,
		})
	}
	return items, total, nil
}
