package core

import (
	"errors"
	"fmt"

	"codebase-app/internal/entity/coreentity"
)

var ErrInvalidSubscriptionTransition = errors.New("invalid subscription transition")

var subscriptionTransitions = map[string]map[string]string{
	"": {
		coreentity.InternalTenantSubscriptionTriggerInitialize: coreentity.InternalTenantSubscriptionStatusFree,
	},
	coreentity.InternalTenantSubscriptionStatusFree: {
		coreentity.InternalTenantSubscriptionTriggerCheckoutCreated: coreentity.InternalTenantSubscriptionStatusPendingActivation,
	},
	coreentity.InternalTenantSubscriptionStatusPendingActivation: {
		coreentity.InternalTenantSubscriptionTriggerInvoicePaidFinal:          coreentity.InternalTenantSubscriptionStatusActive,
		coreentity.InternalTenantSubscriptionTriggerInvoiceExpiredOrAbandoned: coreentity.InternalTenantSubscriptionStatusFree,
	},
	coreentity.InternalTenantSubscriptionStatusActive: {
		coreentity.InternalTenantSubscriptionTriggerRenewalUnpaid:     coreentity.InternalTenantSubscriptionStatusGracePeriod,
		coreentity.InternalTenantSubscriptionTriggerCancelAtPeriodEnd: coreentity.InternalTenantSubscriptionStatusCancelled,
	},
	coreentity.InternalTenantSubscriptionStatusGracePeriod: {
		coreentity.InternalTenantSubscriptionTriggerRenewalPaid: coreentity.InternalTenantSubscriptionStatusActive,
		coreentity.InternalTenantSubscriptionTriggerGraceEnded:  coreentity.InternalTenantSubscriptionStatusSuspended,
	},
	coreentity.InternalTenantSubscriptionStatusSuspended: {
		coreentity.InternalTenantSubscriptionTriggerOutstandingPaidAndReactivated: coreentity.InternalTenantSubscriptionStatusActive,
		coreentity.InternalTenantSubscriptionTriggerAccountClosedWithoutRecovery:  coreentity.InternalTenantSubscriptionStatusExpired,
	},
	coreentity.InternalTenantSubscriptionStatusCancelled: {
		coreentity.InternalTenantSubscriptionTriggerPeriodEnded: coreentity.InternalTenantSubscriptionStatusExpired,
	},
}

func NextSubscriptionStatus(currentStatus string, trigger string) (string, error) {
	nextStatus, ok := subscriptionTransitions[currentStatus][trigger]
	if !ok {
		return "", fmt.Errorf(
			"%w: %s via %s",
			ErrInvalidSubscriptionTransition,
			currentStatus,
			trigger,
		)
	}

	return nextStatus, nil
}

func CanTransitionSubscription(currentStatus string, trigger string) bool {
	_, err := NextSubscriptionStatus(currentStatus, trigger)
	return err == nil
}

func IsTerminalSubscriptionStatus(status string) bool {
	switch status {
	case coreentity.InternalTenantSubscriptionStatusExpired:
		return true
	default:
		return false
	}
}
