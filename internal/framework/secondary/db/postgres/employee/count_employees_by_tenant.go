package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *employeeRepo) CountEmployeesByTenant(ctx context.Context, tenantID string) (int64, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:employee:count_employees_by_tenant:CountEmployeesByTenant",
	)
	defer span.End()

	query := `SELECT COUNT(*) FROM employees WHERE tenant_id = ? AND deleted_at IS NULL`

	exec := r.executor(ctx)
	var count int64
	err := exec.GetContext(ctx, &count, exec.Rebind(query), tenantID)
	return count, err
}
