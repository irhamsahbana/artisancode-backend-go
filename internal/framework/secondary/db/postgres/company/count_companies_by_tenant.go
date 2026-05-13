package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
)

func (r *companyRepo) CountCompaniesByTenant(ctx context.Context, tenantID string) (int64, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:company:count_companies_by_tenant:CountCompaniesByTenant",
	)
	defer span.End()

	query := `SELECT COUNT(*) FROM org_units WHERE tenant_id = ? AND category = 'company' AND deleted_at IS NULL`

	var count int64
	err := r.db.GetContext(ctx, &count, r.db.Rebind(query), tenantID)
	return count, err
}
