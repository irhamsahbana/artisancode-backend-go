package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *userRepo) CountUsersByTenant(ctx context.Context, tenantID string) (int64, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:count_users_by_tenant:CountUsersByTenant",
	)
	defer span.End()

	query := `SELECT COUNT(*) FROM users WHERE tenant_id = ? AND deleted_at IS NULL`

	exec := r.executor(ctx)
	var count int64
	err := exec.GetContext(ctx, &count, exec.Rebind(query), tenantID)
	return count, err
}
