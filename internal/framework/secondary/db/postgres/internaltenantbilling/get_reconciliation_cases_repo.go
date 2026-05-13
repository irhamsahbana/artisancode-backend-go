package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalTenantBillingRepo) GetReconciliationCases(
	ctx context.Context,
	filter coreentity.InternalBillingReconciliationCaseFilter,
) ([]coreentity.InternalBillingReconciliationCase, error) {
	q := `
		SELECT
			id, tenant_id, status, source_type, source_reference_id,
			reason, resolution_note, metadata, resolved_at, created_at, updated_at
		FROM internal_billing_reconciliation_cases
		WHERE 1=1
	`
	args := []any{}

	if filter.Status != "" {
		q += " AND status = ?"
		args = append(args, filter.Status)
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

	var items []coreentity.InternalBillingReconciliationCase
	err := r.db.SelectContext(ctx, &items, q, args...)
	return items, err
}
