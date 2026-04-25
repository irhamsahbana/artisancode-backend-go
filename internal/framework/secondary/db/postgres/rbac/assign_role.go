package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) AssignRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:AssignRole")
	defer span.End()

	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.UserID, data.RoleID)
	return err
}
