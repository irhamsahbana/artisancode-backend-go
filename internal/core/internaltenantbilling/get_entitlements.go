package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetEntitlements(
	ctx context.Context,
) (*coreentity.InternalEntitlementSnapshot, error) {
	_, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_entitlements:GetEntitlements")
	defer span.End()

	userCtx := common.GetUserContext(ctx)

	if c.billingRepo != nil && userCtx.TenantID != "" {
		snap, err := c.billingRepo.GetLatestEntitlementSnapshot(ctx, userCtx.TenantID)
		if err == nil && snap != nil {
			return snap, nil
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

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
