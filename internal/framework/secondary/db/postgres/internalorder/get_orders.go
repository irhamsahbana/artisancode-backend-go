package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalOrderRepo) GetOrders(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalOrder, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalorder:get_orders:GetOrders")
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	where := `WHERE deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `SELECT COUNT(*) FROM internal_orders ` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]orderRow, 0)
	query := `
		SELECT
			id,
			tenant_id,
			company_id,
			order_number,
			status,
			currency_code,
			subtotal_amount,
			discount_amount,
			tax_amount,
			total_amount,
			internal_product_id,
			internal_product_pricing_id,
			pricing_snapshot,
			source_type,
			source_reference_id,
			metadata,
			ordered_at,
			created_at,
			updated_at
		FROM internal_orders
		` + where + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalOrder, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapOrder(item))
	}
	return items, total, nil
}
