package repository

import (
	"context"
	"fmt"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalInvoiceRepo) getPaymentAttemptsByWhere(
	ctx context.Context,
	tenantID string,
	where string,
	value string,
	limit int,
) ([]coreentity.InternalPaymentAttempt, error) {
	rows := make([]paymentAttemptRow, 0)
	args := []any{value}
	query := fmt.Sprintf(`
		SELECT
			pa.id,
			pa.internal_invoice_id,
			pa.provider,
			pa.payment_method_type,
			pa.payment_channel_code,
			pa.provider_reference,
			pa.provider_request_id,
			pa.provider_payment_url,
			pa.provider_payload_snapshot,
			pa.status,
			pa.requested_amount,
			pa.paid_amount,
			pa.expired_at,
			pa.paid_at,
			pa.failed_at,
			pa.raw_last_status,
			pa.metadata,
			pa.created_at,
			pa.updated_at
		FROM internal_payment_attempts pa
		JOIN internal_invoices i
			ON i.id = pa.internal_invoice_id
			AND i.deleted_at IS NULL
		JOIN internal_orders o
			ON o.id = i.internal_order_id
			AND o.deleted_at IS NULL
		WHERE %s
			AND pa.deleted_at IS NULL
	`, where)
	if tenantID != "" {
		query += ` AND o.tenant_id = ?`
		args = append(args, tenantID)
	}
	query += ` ORDER BY pa.created_at DESC LIMIT ?`
	args = append(args, limit)
	if err := r.exec(ctx).SelectContext(ctx, &rows, r.exec(ctx).Rebind(query), args...); err != nil {
		return nil, err
	}
	items := make([]coreentity.InternalPaymentAttempt, 0, len(rows))
	for _, row := range rows {
		items = append(items, *mapPaymentAttempt(row))
	}
	return items, nil
}
