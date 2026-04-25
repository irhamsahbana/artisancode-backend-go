package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *rbacRepo) CreateRole(ctx context.Context, data coreentity.Role) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:CreateRole")
	defer span.End()

	query := `
		INSERT INTO roles (tenant_id, name)
		VALUES (?, ?)
		RETURNING id, tenant_id, name
	`

	var item coreentity.Role
	if err := r.db.GetContext(ctx, &item, r.db.Rebind(query), data.TenantID, data.Name); err != nil {
		return nil, err
	}
	return &item, nil
}
