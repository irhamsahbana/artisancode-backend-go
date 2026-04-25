package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) RemovePermission(ctx context.Context, userID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:RemovePermission")
	defer span.End()

	query := `
		DELETE FROM user_permissions
		WHERE user_id = ? AND permission_id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), userID, permissionID)
	return err
}
