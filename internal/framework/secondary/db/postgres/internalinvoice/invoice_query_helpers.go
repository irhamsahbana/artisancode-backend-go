package repository

import (
	"context"
	"database/sql"
	"fmt"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"
)

func (r *internalInvoiceRepo) getInvoiceByWhere(
	ctx context.Context,
	tenantID string,
	where string,
	value string,
) (*coreentity.InternalInvoice, error) {
	var item invoiceRow
	query := fmt.Sprintf(`
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
		JOIN internal_orders o
			ON o.id = i.internal_order_id
			AND o.deleted_at IS NULL
		WHERE %s
			AND o.tenant_id = ?
			AND i.deleted_at IS NULL
	`, where)
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), value, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Invoice not found")
		}
		return nil, err
	}
	return mapInvoice(item), nil
}
