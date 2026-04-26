package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) MarkInvoicePaid(
	ctx context.Context,
	tenantID string,
	id string,
	amountPaid string,
) (*coreentity.InternalInvoice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:mark_invoice_paid:MarkInvoicePaid",
	)
	defer span.End()

	query := `
		UPDATE internal_invoices i
		SET
			status = 'paid',
			amount_paid = ?,
			amount_outstanding = 0,
			paid_at = NOW(),
			updated_at = NOW()
		FROM internal_orders o
		WHERE i.internal_order_id = o.id
			AND i.id = ?
			AND o.tenant_id = ?
			AND i.deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), amountPaid, id, tenantID); err != nil {
		return nil, err
	}
	return r.GetInvoice(ctx, tenantID, id)
}
