package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) SetRolePermissions(ctx context.Context, roleID, tenantID string, permissionIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:SetRolePermissions")
	defer span.End()

	// Delete existing permissions
	deleteQuery := `DELETE FROM role_permissions WHERE role_id = ?`
	if _, err := r.db.ExecContext(ctx, r.db.Rebind(deleteQuery), roleID); err != nil {
		return err
	}

	// Insert new permissions
	if len(permissionIDs) > 0 {
		insertQuery := `INSERT INTO role_permissions (role_id, permission_id, tenant_id) VALUES (?, ?, ?)`
		for _, permID := range permissionIDs {
			if _, err := r.db.ExecContext(ctx, r.db.Rebind(insertQuery), roleID, permID, tenantID); err != nil {
				return err
			}
		}
	}
	return nil
}
