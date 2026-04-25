package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) GetUserPermissions(ctx context.Context, userID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetUserPermissions")
	defer span.End()

	query := `
		SELECT DISTINCT p.id, p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		INNER JOIN roles r ON r.id = rp.role_id
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL
		ORDER BY p.name ASC
	`

	rows := make([]coreentity.Permission, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), userID); err != nil {
		return nil, err
	}
	return rows, nil
}
