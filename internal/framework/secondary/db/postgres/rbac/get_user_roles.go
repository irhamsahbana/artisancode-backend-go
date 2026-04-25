package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) GetUserRoles(ctx context.Context, filter coreentity.UserRoleFilter) ([]coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetUserRoles")
	defer span.End()

	query := `
		SELECT r.id, r.tenant_id, r.name
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL
		ORDER BY r.name ASC
	`

	rows := make([]coreentity.Role, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), filter.UserID); err != nil {
		return nil, err
	}
	return rows, nil
}
