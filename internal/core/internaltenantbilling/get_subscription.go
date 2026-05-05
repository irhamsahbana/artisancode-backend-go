package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetSubscription(
	ctx context.Context,
) (*coreentity.TenantBillingSubscriptionView, error) {
	_, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_subscription:GetSubscription")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	if userCtx.TenantID == "" {
		return &coreentity.TenantBillingSubscriptionView{
			Status: coreentity.InternalTenantSubscriptionStatusFree,
		}, nil
	}

	if c.billingRepo == nil {
		return &coreentity.TenantBillingSubscriptionView{
			Status: coreentity.InternalTenantSubscriptionStatusFree,
		}, nil
	}

	sub, err := c.billingRepo.GetSubscription(ctx, userCtx.TenantID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if sub == nil || errors.Is(err, sql.ErrNoRows) {
		return &coreentity.TenantBillingSubscriptionView{
			Status: coreentity.InternalTenantSubscriptionStatusFree,
		}, nil
	}

	planName := ""
	billingCycle := ""
	if sub.PlanSnapshot != nil {
		if name, ok := sub.PlanSnapshot["name"].(string); ok {
			planName = name
		}
		if cycle, ok := sub.PlanSnapshot["billing_cycle"].(string); ok {
			billingCycle = cycle
		}
	}

	return &coreentity.TenantBillingSubscriptionView{
		ID:                 sub.ID,
		Status:             sub.Status,
		PlanName:           planName,
		BillingCycle:       billingCycle,
		CurrentPeriodStart: sub.CurrentPeriodStartedAt,
		CurrentPeriodEnd:   sub.CurrentPeriodEndedAt,
		NextRenewalAt:      sub.CurrentPeriodEndedAt,
	}, nil
}
