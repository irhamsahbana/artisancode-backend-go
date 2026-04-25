package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) UpdateRole(ctx context.Context, roleID, tenantID string, permissionIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:UpdateRole")
	defer span.End()

	if _, err := r.GetRoleWithPermissions(ctx, roleID, tenantID); err != nil {
		return err
	}

	return r.SetRolePermissions(ctx, roleID, tenantID, permissionIDs)
}
