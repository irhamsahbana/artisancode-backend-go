package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) ExecuteSubscriptionAction(
	ctx context.Context,
	input coreentity.TenantBillingSubscriptionActionInput,
) (*coreentity.TenantBillingSubscriptionActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:execute_subscription_action:ExecuteSubscriptionAction")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	input.UserCtx = userCtx

	sub, err := c.billingRepo.GetSubscription(ctx, userCtx.TenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageSubscriptionNotFound))
		}
		return nil, err
	}

	switch input.Action {
	case coreentity.TenantBillingSubscriptionActionCancel:
		return c.cancelSubscription(ctx, sub)
	case coreentity.TenantBillingSubscriptionActionReactivate:
		return c.reactivateSubscription(ctx, sub)
	default:
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedSubscriptionAction))
	}
}

func (c *internalTenantBillingCore) cancelSubscription(
	ctx context.Context,
	sub *coreentity.InternalTenantSubscription,
) (*coreentity.TenantBillingSubscriptionActionResult, error) {
	if sub.Status != coreentity.InternalTenantSubscriptionStatusActive &&
		sub.Status != coreentity.InternalTenantSubscriptionStatusGracePeriod {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedSubscriptionAction))
	}

	if err := c.billingRepo.SetSubscriptionStatus(ctx, sub.TenantID, sub.ID, coreentity.InternalTenantSubscriptionStatusCancelled); err != nil {
		return nil, err
	}

	return &coreentity.TenantBillingSubscriptionActionResult{
		SubscriptionID: sub.ID,
		Status:         coreentity.InternalTenantSubscriptionStatusCancelled,
	}, nil
}

func (c *internalTenantBillingCore) reactivateSubscription(
	ctx context.Context,
	sub *coreentity.InternalTenantSubscription,
) (*coreentity.TenantBillingSubscriptionActionResult, error) {
	if sub.Status != coreentity.InternalTenantSubscriptionStatusSuspended {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedSubscriptionAction))
	}

	if err := c.billingRepo.SetSubscriptionStatus(ctx, sub.TenantID, sub.ID, coreentity.InternalTenantSubscriptionStatusActive); err != nil {
		return nil, err
	}

	return &coreentity.TenantBillingSubscriptionActionResult{
		SubscriptionID: sub.ID,
		Status:         coreentity.InternalTenantSubscriptionStatusActive,
	}, nil
}
