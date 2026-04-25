package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetPermissionsByRoleID")
	defer span.End()

	query := `
		SELECT p.id, p.name, COALESCE(p.description, '') AS description
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ? AND p.deleted_at IS NULL
		ORDER BY p.name ASC
	`

	rows := make([]coreentity.Permission, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), roleID); err != nil {
		return nil, err
	}
	return rows, nil
}
