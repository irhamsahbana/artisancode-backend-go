package core

import (
	"errors"
	"fmt"
	"strings"

	"codebase-app/internal/entity/coreentity"
)

var ErrDOKUStatusUnmapped = errors.New("doku status unmapped")

func NormalizeDOKUPaymentStatus(
	providerStatus string,
	providerContext string,
) (coreentity.InternalBillingPaymentStatusNormalization, error) {
	normalized := coreentity.InternalBillingPaymentStatusNormalization{
		ProviderStatus:  strings.TrimSpace(providerStatus),
		ProviderContext: strings.TrimSpace(providerContext),
	}

	switch normalized.ProviderContext {
	case coreentity.DOKUProviderContextCheckoutOrder:
		return normalizeDOKUCheckoutOrderStatus(normalized)
	case coreentity.DOKUProviderContextTransactionChannel:
		return normalizeDOKUTransactionChannelStatus(normalized)
	case coreentity.DOKUProviderContextOrderCheck:
		return normalizeDOKUOrderCheckStatus(normalized)
	default:
		return coreentity.InternalBillingPaymentStatusNormalization{}, fmt.Errorf(
			"%w: %s/%s",
			ErrDOKUStatusUnmapped,
			providerContext,
			providerStatus,
		)
	}
}

func normalizeDOKUCheckoutOrderStatus(
	normalized coreentity.InternalBillingPaymentStatusNormalization,
) (coreentity.InternalBillingPaymentStatusNormalization, error) {
	switch normalized.ProviderStatus {
	case coreentity.DOKUProviderStatusSuccess:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusSucceeded
		normalized.IsTerminal = true
	case coreentity.DOKUProviderStatusExpired:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusExpired
		normalized.IsTerminal = true
	case coreentity.DOKUProviderStatusCancelled:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusCancelled
		normalized.PaymentEventStatus = coreentity.InternalBillingPaymentEventStatusCancelledByMerchant
		normalized.IsTerminal = true
	case coreentity.DOKUProviderStatusPending:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusPending
	default:
		return coreentity.InternalBillingPaymentStatusNormalization{}, fmt.Errorf(
			"%w: %s/%s",
			ErrDOKUStatusUnmapped,
			normalized.ProviderContext,
			normalized.ProviderStatus,
		)
	}

	return normalized, nil
}

func normalizeDOKUTransactionChannelStatus(
	normalized coreentity.InternalBillingPaymentStatusNormalization,
) (coreentity.InternalBillingPaymentStatusNormalization, error) {
	switch normalized.ProviderStatus {
	case coreentity.DOKUProviderStatusTransactionOK:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusSucceeded
		normalized.IsTerminal = true
	case coreentity.DOKUProviderStatusTransactionNOK:
		normalized.PaymentEventStatus = coreentity.InternalBillingPaymentEventStatusFailedByChannel
	default:
		return coreentity.InternalBillingPaymentStatusNormalization{}, fmt.Errorf(
			"%w: %s/%s",
			ErrDOKUStatusUnmapped,
			normalized.ProviderContext,
			normalized.ProviderStatus,
		)
	}

	return normalized, nil
}

func normalizeDOKUOrderCheckStatus(
	normalized coreentity.InternalBillingPaymentStatusNormalization,
) (coreentity.InternalBillingPaymentStatusNormalization, error) {
	switch normalized.ProviderStatus {
	case coreentity.DOKUProviderStatusOrderGenerated:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusPending
	case coreentity.DOKUProviderStatusOrderExpired:
		normalized.PaymentAttemptStatus = coreentity.PaymentAttemptStatusExpired
		normalized.IsTerminal = true
	case coreentity.DOKUProviderStatusOrderRecovered:
		normalized.PaymentEventStatus = coreentity.InternalBillingPaymentEventStatusRecovered
	default:
		return coreentity.InternalBillingPaymentStatusNormalization{}, fmt.Errorf(
			"%w: %s/%s",
			ErrDOKUStatusUnmapped,
			normalized.ProviderContext,
			normalized.ProviderStatus,
		)
	}

	return normalized, nil
}
