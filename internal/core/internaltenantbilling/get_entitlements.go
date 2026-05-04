package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetEntitlements(
	ctx context.Context,
) (*coreentity.InternalEntitlementSnapshot, error) {
	_, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_entitlements:GetEntitlements")
	defer span.End()

	snapshot, err := CalculateEntitlementSnapshot(coreentity.InternalEntitlementCalculationInput{
		SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusFree,
		FreeTier: coreentity.InternalEntitlementGrant{
			Features: []string{"attendance_basic"},
			UsageLimits: map[string]int64{
				"employees": 10,
				"branches":  1,
			},
		},
		SafeAccess: coreentity.InternalEntitlementGrant{
			Features: []string{"account_recovery"},
			UsageLimits: map[string]int64{
				"employees": 0,
				"branches":  0,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}
