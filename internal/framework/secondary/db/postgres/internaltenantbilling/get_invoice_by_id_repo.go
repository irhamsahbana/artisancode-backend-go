package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalTenantBillingRepo) GetInvoiceByID(
	ctx context.Context,
	id string,
) (*coreentity.InternalTenantInvoice, error) {
	q := `
		SELECT
			id, tenant_id, internal_billing_account_id, internal_tenant_subscription_id,
			internal_tenant_subscription_change_id, invoice_number, status, currency_code,
			amount, amount_paid, amount_outstanding, source_type, source_reference_id,
			target_subscription_state, due_at, paid_at, expired_at, metadata,
			created_at, updated_at
		FROM internal_tenant_invoices
		WHERE id = ?
	`
	var item coreentity.InternalTenantInvoice
	err := r.db.GetContext(ctx, &item, q, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *internalTenantBillingRepo) CreateInternalInvoice(
	ctx context.Context,
	data coreentity.InternalTenantInvoice,
) (*coreentity.InternalTenantInvoice, error) {
	q := `
		INSERT INTO internal_tenant_invoices (
			tenant_id, internal_billing_account_id, internal_tenant_subscription_id,
			invoice_number, status, currency_code, amount, amount_paid, amount_outstanding,
			source_type, source_reference_id, target_subscription_state, metadata
		) VALUES (
			:tenant_id, :internal_billing_account_id, :internal_tenant_subscription_id,
			:invoice_number, :status, :currency_code, :amount, :amount_paid, :amount_outstanding,
			:source_type, :source_reference_id, :target_subscription_state, :metadata
		)
		RETURNING
			id, tenant_id, internal_billing_account_id, internal_tenant_subscription_id,
			internal_tenant_subscription_change_id, invoice_number, status, currency_code,
			amount, amount_paid, amount_outstanding, source_type, source_reference_id,
			target_subscription_state, due_at, paid_at, expired_at, metadata,
			created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var item coreentity.InternalTenantInvoice
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, rows.Err()
}
