package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalQuotationRepo) GetQuotations(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalQuotation, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:get_quotations:GetQuotations",
	)
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	where := `WHERE deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `SELECT COUNT(*) FROM internal_quotations ` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]quotationRow, 0)
	query := `
		SELECT
			id,
			tenant_id,
			company_id,
			quotation_number,
			status,
			currency_code,
			subtotal_amount,
			discount_amount,
			tax_amount,
			total_amount,
			internal_product_id,
			internal_product_pricing_id,
			pricing_snapshot,
			quote_snapshot,
			expires_at,
			approved_at,
			converted_to_order_id,
			metadata,
			created_at,
			updated_at
		FROM internal_quotations
		` + where + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalQuotation, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapQuotation(item))
	}
	return items, total, nil
}
