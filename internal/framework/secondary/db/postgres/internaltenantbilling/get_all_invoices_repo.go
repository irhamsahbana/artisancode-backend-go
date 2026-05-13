package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalTenantBillingRepo) GetAllInvoices(
	ctx context.Context,
	filter coreentity.InternalBillingInvoiceListFilter,
) ([]coreentity.InternalTenantInvoice, error) {
	q := `
		SELECT
			id, tenant_id, internal_billing_account_id, internal_tenant_subscription_id,
			internal_tenant_subscription_change_id, invoice_number, status, currency_code,
			amount, amount_paid, amount_outstanding, source_type, source_reference_id,
			target_subscription_state, due_at, paid_at, expired_at, metadata,
			created_at, updated_at
		FROM internal_tenant_invoices
		WHERE 1=1
	`
	args := []any{}

	if filter.TenantID != nil && *filter.TenantID != "" {
		q += " AND tenant_id = ?"
		args = append(args, *filter.TenantID)
	}
	if filter.Status != "" {
		q += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.FromDate != "" {
		q += " AND created_at >= ?"
		args = append(args, filter.FromDate)
	}
	if filter.ToDate != "" {
		q += " AND created_at <= ?"
		args = append(args, filter.ToDate)
	}

	if filter.Size <= 0 {
		filter.Size = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Size
	q += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Size, offset)

	var items []coreentity.InternalTenantInvoice
	err := r.db.SelectContext(ctx, &items, q, args...)
	return items, err
}
