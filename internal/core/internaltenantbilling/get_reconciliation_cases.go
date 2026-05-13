package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetReconciliationCases(
	ctx context.Context,
	filter coreentity.InternalBillingReconciliationCaseFilter,
) ([]coreentity.InternalBillingReconciliationCase, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_reconciliation_cases:GetReconciliationCases")
	defer span.End()

	return c.billingRepo.GetReconciliationCases(ctx, filter)
}
