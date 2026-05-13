package repository

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
)

func (r *internalTenantBillingRepo) GetLedgerEntries(
	ctx context.Context,
	tenantID string,
	filter coreentity.InternalBillingLedgerListFilter,
) ([]coreentity.InternalBillingLedgerEntry, error) {
	q := `
		SELECT
			id, tenant_id, internal_tenant_subscription_id,
			entry_type, source_type, source_reference_id,
			occurred_at, metadata, created_at, updated_at
		FROM internal_billing_ledger_entries
		WHERE tenant_id = ?
	`
	args := []any{tenantID}

	if len(filter.EntryTypes) > 0 {
		placeholders := make([]string, len(filter.EntryTypes))
		for i, et := range filter.EntryTypes {
			placeholders[i] = "?"
			args = append(args, et)
		}
		q += " AND entry_type IN (" + strings.Join(placeholders, ",") + ")"
	}
	if filter.FromDate != "" {
		q += " AND occurred_at >= ?"
		args = append(args, filter.FromDate)
	}
	if filter.ToDate != "" {
		q += " AND occurred_at <= ?"
		args = append(args, filter.ToDate)
	}

	if filter.Size <= 0 {
		filter.Size = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Size
	q += " ORDER BY occurred_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Size, offset)

	var items []coreentity.InternalBillingLedgerEntry
	err := r.db.SelectContext(ctx, &items, q, args...)
	return items, err
}
