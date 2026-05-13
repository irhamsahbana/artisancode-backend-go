package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetLedgerEntries(
	ctx context.Context,
	tenantID string,
	filter coreentity.InternalBillingLedgerListFilter,
) ([]coreentity.InternalBillingLedgerEntry, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_ledger_entries:GetLedgerEntries")
	defer span.End()

	return c.billingRepo.GetLedgerEntries(ctx, tenantID, filter)
}
