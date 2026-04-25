package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *rbacRepo) GetRoleWithPermissions(ctx context.Context, roleID, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetRoleWithPermissions")
	defer span.End()

	query := `
		SELECT id, tenant_id, name
		FROM roles
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var item coreentity.Role
	if err := r.db.GetContext(ctx, &item, r.db.Rebind(query), roleID, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Role not found")
		}
		return nil, err
	}
	return &item, nil
}
