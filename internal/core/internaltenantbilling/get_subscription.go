package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetSubscription(
	ctx context.Context,
) (*coreentity.TenantBillingSubscriptionView, error) {
	_, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_subscription:GetSubscription")
	defer span.End()

	return &coreentity.TenantBillingSubscriptionView{
		Status: coreentity.InternalTenantSubscriptionStatusFree,
	}, nil
}
