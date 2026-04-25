package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) AssignPermission(ctx context.Context, userID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:AssignPermission")
	defer span.End()

	query := `
		INSERT INTO user_permissions (user_id, permission_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, permission_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), userID, permissionID)
	return err
}
