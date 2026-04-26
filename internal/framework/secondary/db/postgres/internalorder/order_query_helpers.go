package repository

import (
	"context"
	"database/sql"
	"fmt"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"
)

func (r *internalOrderRepo) getOrderByWhere(
	ctx context.Context,
	tenantID string,
	where string,
	value string,
) (*coreentity.InternalOrder, error) {
	var item orderRow
	query := fmt.Sprintf(`
		SELECT
			o.id,
			o.tenant_id,
			o.company_id,
			o.order_number,
			o.status,
			o.currency_code,
			o.subtotal_amount,
			o.discount_amount,
			o.tax_amount,
			o.total_amount,
			o.internal_product_id,
			o.internal_product_pricing_id,
			o.pricing_snapshot,
			o.source_type,
			o.source_reference_id,
			o.metadata,
			o.ordered_at,
			o.created_at,
			o.updated_at
		FROM internal_orders o
		LEFT JOIN internal_invoices i
			ON i.internal_order_id = o.id
			AND i.deleted_at IS NULL
		WHERE %s
			AND o.tenant_id = ?
			AND o.deleted_at IS NULL
	`, where)
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), value, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Order not found")
		}
		return nil, err
	}
	return mapOrder(item), nil
}
