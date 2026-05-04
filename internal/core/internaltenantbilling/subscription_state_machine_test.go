package core

import (
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

func TestNextSubscriptionStatus(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		trigger       string
		want          string
	}{
		{
			name:    "initializes free status",
			trigger: coreentity.InternalTenantSubscriptionTriggerInitialize,
			want:    coreentity.InternalTenantSubscriptionStatusFree,
		},
		{
			name:          "moves free checkout to pending activation",
			currentStatus: coreentity.InternalTenantSubscriptionStatusFree,
			trigger:       coreentity.InternalTenantSubscriptionTriggerCheckoutCreated,
			want:          coreentity.InternalTenantSubscriptionStatusPendingActivation,
		},
		{
			name:          "activates after final paid invoice",
			currentStatus: coreentity.InternalTenantSubscriptionStatusPendingActivation,
			trigger:       coreentity.InternalTenantSubscriptionTriggerInvoicePaidFinal,
			want:          coreentity.InternalTenantSubscriptionStatusActive,
		},
		{
			name:          "falls back to free after abandoned checkout",
			currentStatus: coreentity.InternalTenantSubscriptionStatusPendingActivation,
			trigger:       coreentity.InternalTenantSubscriptionTriggerInvoiceExpiredOrAbandoned,
			want:          coreentity.InternalTenantSubscriptionStatusFree,
		},
		{
			name:          "starts grace after unpaid renewal",
			currentStatus: coreentity.InternalTenantSubscriptionStatusActive,
			trigger:       coreentity.InternalTenantSubscriptionTriggerRenewalUnpaid,
			want:          coreentity.InternalTenantSubscriptionStatusGracePeriod,
		},
		{
			name:          "recovers from grace after renewal payment",
			currentStatus: coreentity.InternalTenantSubscriptionStatusGracePeriod,
			trigger:       coreentity.InternalTenantSubscriptionTriggerRenewalPaid,
			want:          coreentity.InternalTenantSubscriptionStatusActive,
		},
		{
			name:          "suspends after grace ended",
			currentStatus: coreentity.InternalTenantSubscriptionStatusGracePeriod,
			trigger:       coreentity.InternalTenantSubscriptionTriggerGraceEnded,
			want:          coreentity.InternalTenantSubscriptionStatusSuspended,
		},
		{
			name:          "reactivates suspended after outstanding paid",
			currentStatus: coreentity.InternalTenantSubscriptionStatusSuspended,
			trigger:       coreentity.InternalTenantSubscriptionTriggerOutstandingPaidAndReactivated,
			want:          coreentity.InternalTenantSubscriptionStatusActive,
		},
		{
			name:          "expires cancelled after period ended",
			currentStatus: coreentity.InternalTenantSubscriptionStatusCancelled,
			trigger:       coreentity.InternalTenantSubscriptionTriggerPeriodEnded,
			want:          coreentity.InternalTenantSubscriptionStatusExpired,
		},
		{
			name:          "expires suspended after closure",
			currentStatus: coreentity.InternalTenantSubscriptionStatusSuspended,
			trigger:       coreentity.InternalTenantSubscriptionTriggerAccountClosedWithoutRecovery,
			want:          coreentity.InternalTenantSubscriptionStatusExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextSubscriptionStatus(tt.currentStatus, tt.trigger)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNextSubscriptionStatusRejectsInvalidTransition(t *testing.T) {
	got, err := NextSubscriptionStatus(
		coreentity.InternalTenantSubscriptionStatusFree,
		coreentity.InternalTenantSubscriptionTriggerInvoicePaidFinal,
	)

	require.Empty(t, got)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidSubscriptionTransition))
	require.False(t, CanTransitionSubscription(
		coreentity.InternalTenantSubscriptionStatusFree,
		coreentity.InternalTenantSubscriptionTriggerInvoicePaidFinal,
	))
}

func TestIsTerminalSubscriptionStatus(t *testing.T) {
	require.True(t, IsTerminalSubscriptionStatus(coreentity.InternalTenantSubscriptionStatusExpired))
	require.False(t, IsTerminalSubscriptionStatus(coreentity.InternalTenantSubscriptionStatusSuspended))
}
