package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *rbacRepo) RemoveRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:RemoveRole")
	defer span.End()

	query := `
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.UserID, data.RoleID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("User role not found")
	}
	return nil
}
