package core

import (
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

func TestNormalizeDOKUPaymentStatus(t *testing.T) {
	tests := []struct {
		name                 string
		providerStatus       string
		providerContext      string
		paymentAttemptStatus string
		paymentEventStatus   string
		isTerminal           bool
	}{
		{
			name:                 "checkout success is succeeded terminal",
			providerStatus:       coreentity.DOKUProviderStatusSuccess,
			providerContext:      coreentity.DOKUProviderContextCheckoutOrder,
			paymentAttemptStatus: coreentity.PaymentAttemptStatusSucceeded,
			isTerminal:           true,
		},
		{
			name:                 "checkout pending remains non-terminal",
			providerStatus:       coreentity.DOKUProviderStatusPending,
			providerContext:      coreentity.DOKUProviderContextCheckoutOrder,
			paymentAttemptStatus: coreentity.PaymentAttemptStatusPending,
		},
		{
			name:                 "checkout cancelled is terminal",
			providerStatus:       coreentity.DOKUProviderStatusCancelled,
			providerContext:      coreentity.DOKUProviderContextCheckoutOrder,
			paymentAttemptStatus: coreentity.PaymentAttemptStatusCancelled,
			paymentEventStatus:   coreentity.InternalBillingPaymentEventStatusCancelledByMerchant,
			isTerminal:           true,
		},
		{
			name:               "transaction failed is channel event only",
			providerStatus:     coreentity.DOKUProviderStatusTransactionNOK,
			providerContext:    coreentity.DOKUProviderContextTransactionChannel,
			paymentEventStatus: coreentity.InternalBillingPaymentEventStatusFailedByChannel,
		},
		{
			name:                 "order generated keeps pending payable flow",
			providerStatus:       coreentity.DOKUProviderStatusOrderGenerated,
			providerContext:      coreentity.DOKUProviderContextOrderCheck,
			paymentAttemptStatus: coreentity.PaymentAttemptStatusPending,
		},
		{
			name:                 "order expired is expired terminal",
			providerStatus:       coreentity.DOKUProviderStatusOrderExpired,
			providerContext:      coreentity.DOKUProviderContextOrderCheck,
			paymentAttemptStatus: coreentity.PaymentAttemptStatusExpired,
			isTerminal:           true,
		},
		{
			name:               "order recovered is recovery event only",
			providerStatus:     coreentity.DOKUProviderStatusOrderRecovered,
			providerContext:    coreentity.DOKUProviderContextOrderCheck,
			paymentEventStatus: coreentity.InternalBillingPaymentEventStatusRecovered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeDOKUPaymentStatus(tt.providerStatus, tt.providerContext)

			require.NoError(t, err)
			require.Equal(t, tt.providerStatus, got.ProviderStatus)
			require.Equal(t, tt.providerContext, got.ProviderContext)
			require.Equal(t, tt.paymentAttemptStatus, got.PaymentAttemptStatus)
			require.Equal(t, tt.paymentEventStatus, got.PaymentEventStatus)
			require.Equal(t, tt.isTerminal, got.IsTerminal)
		})
	}
}

func TestNormalizeDOKUPaymentStatusRejectsUnmappedStatus(t *testing.T) {
	got, err := NormalizeDOKUPaymentStatus(
		"NEW_STATUS",
		coreentity.DOKUProviderContextCheckoutOrder,
	)

	require.Empty(t, got)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrDOKUStatusUnmapped))
}
