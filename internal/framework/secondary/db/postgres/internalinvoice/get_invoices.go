package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) GetInvoices(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalInvoice, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_invoices:GetInvoices",
	)
	defer span.End()

	filter = normalizeCommerceListFilter(filter)
	where := `WHERE i.deleted_at IS NULL AND o.deleted_at IS NULL`
	args := []any{}
	if filter.TenantID != "" {
		where += ` AND o.tenant_id = ?`
		args = append(args, filter.TenantID)
	}
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM internal_invoices i
		JOIN internal_orders o ON o.id = i.internal_order_id
		` + where
	if err := r.exec(ctx).GetContext(ctx, &total, r.exec(ctx).Rebind(countQuery), args...); err != nil {
		return nil, 0, err
	}
	rows := make([]invoiceRow, 0)
	query := `
		SELECT
			i.id,
			i.internal_order_id,
			i.invoice_number,
			i.status,
			i.currency_code,
			i.amount,
			i.amount_paid,
			i.amount_outstanding,
			i.due_at,
			i.paid_at,
			i.expired_at,
			i.metadata,
			i.created_at,
			i.updated_at
		FROM internal_invoices i
		JOIN internal_orders o ON o.id = i.internal_order_id
		` + where + `
		ORDER BY i.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, 0, err
	}
	items := make([]coreentity.InternalInvoice, 0, len(rows))
	for _, item := range rows {
		items = append(items, *mapInvoice(item))
	}
	return items, total, nil
}
