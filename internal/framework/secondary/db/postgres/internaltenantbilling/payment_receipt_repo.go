package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalTenantBillingRepo) CreatePaymentReceipt(
	ctx context.Context,
	data coreentity.InternalPaymentReceipt,
) (*coreentity.InternalPaymentReceipt, error) {
	q := `
		INSERT INTO internal_payment_receipts (
			internal_invoice_id, internal_payment_attempt_id, receipt_number,
			status, amount_received, currency_code, received_at, source_type,
			reference_number, notes, metadata
		) VALUES (
			:internal_invoice_id, :internal_payment_attempt_id, :receipt_number,
			:status, :amount_received, :currency_code, :received_at, :source_type,
			:reference_number, :notes, :metadata
		)
		RETURNING
			id, internal_invoice_id, internal_payment_attempt_id, receipt_number,
			status, amount_received, currency_code, received_at, verified_at,
			verified_by_user_id, source_type, reference_number, notes, metadata,
			created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var item coreentity.InternalPaymentReceipt
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, rows.Err()
}

func (r *internalTenantBillingRepo) GetPaymentReceipt(
	ctx context.Context,
	id string,
) (*coreentity.InternalPaymentReceipt, error) {
	q := `
		SELECT
			id, internal_invoice_id, internal_payment_attempt_id, receipt_number,
			status, amount_received, currency_code, received_at, verified_at,
			verified_by_user_id, source_type, reference_number, notes, metadata,
			created_at, updated_at
		FROM internal_payment_receipts
		WHERE id = ?
	`
	var item coreentity.InternalPaymentReceipt
	err := r.db.GetContext(ctx, &item, q, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *internalTenantBillingRepo) UpdatePaymentReceipt(
	ctx context.Context,
	data coreentity.InternalPaymentReceipt,
) (*coreentity.InternalPaymentReceipt, error) {
	q := `
		UPDATE internal_payment_receipts
		SET status = :status, verified_at = :verified_at, verified_by_user_id = :verified_by_user_id
		WHERE id = :id
		RETURNING
			id, internal_invoice_id, internal_payment_attempt_id, receipt_number,
			status, amount_received, currency_code, received_at, verified_at,
			verified_by_user_id, source_type, reference_number, notes, metadata,
			created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var item coreentity.InternalPaymentReceipt
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, rows.Err()
}
