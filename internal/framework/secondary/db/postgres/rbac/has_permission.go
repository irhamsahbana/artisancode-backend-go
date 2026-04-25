package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:HasPermission")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM role_permissions rp
			INNER JOIN roles r ON r.id = rp.role_id
			INNER JOIN permissions p ON p.id = rp.permission_id
			INNER JOIN user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = ? AND p.name = ? AND r.deleted_at IS NULL
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), userID, permission); err != nil {
		return false, err
	}
	return exists, nil
}
