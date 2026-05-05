package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	billingCore "codebase-app/internal/core/internaltenantbilling"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (c *webhookCore) HandleDOKUWebhook(
	ctx context.Context,
	notification coreentity.DOKUWebhookNotification,
) (*coreentity.DOKUWebhookResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:webhook:handle_doku_webhook:HandleDOKUWebhook")
	defer span.End()

	if c.dokuVerifier == nil {
		err := errors.New("doku webhook verifier is required")
		tracing.RecordError(span, err)
		return nil, err
	}

	if !c.dokuVerifier.VerifyWebhookSignatureHeaders(
		notification.Headers,
		notification.RawBody,
		notification.TargetPath,
	) {
		err := errors.New("invalid doku webhook signature")
		log.Ctx(ctx).Warn().
			Str("invoice_number", notification.Event.Order.InvoiceNumber).
			Str("request_id", notification.Headers.RequestID).
			Msg("DOKU webhook signature verification failed")
		tracing.RecordError(span, err)
		return nil, err
	}

	log.Ctx(ctx).Info().
		Str("invoice_number", notification.Event.Order.InvoiceNumber).
		Str("order_status", notification.Event.Order.Status).
		Str("payment_status", notification.Event.Transaction.Status).
		Str("request_id", notification.Headers.RequestID).
		Msg("DOKU webhook accepted")

	normalized, err := billingCore.NormalizeDOKUPaymentStatus(
		notification.Event.Transaction.Status,
		coreentity.DOKUProviderContextCheckoutOrder,
	)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).
			Str("invoice_number", notification.Event.Order.InvoiceNumber).
			Str("raw_status", notification.Event.Transaction.Status).
			Msg("DOKU webhook status unmapped, acknowledged without processing")
		return &coreentity.DOKUWebhookResult{
			Provider:      "doku",
			InvoiceNumber: notification.Event.Order.InvoiceNumber,
			OrderStatus:   notification.Event.Order.Status,
			PaymentStatus: notification.Event.Transaction.Status,
		}, nil
	}

	if c.billingRepo == nil {
		return &coreentity.DOKUWebhookResult{
			Provider:      "doku",
			InvoiceNumber: notification.Event.Order.InvoiceNumber,
			OrderStatus:   notification.Event.Order.Status,
			PaymentStatus: notification.Event.Transaction.Status,
		}, nil
	}

	attempt, err := c.billingRepo.GetPaymentAttemptByProviderRef(ctx, notification.Event.Order.InvoiceNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup payment attempt by %s: %w", notification.Event.Order.InvoiceNumber, err)
	}

	rawLastStatusJSON, _ := json.Marshal(map[string]any{
		"order_status":   notification.Event.Order.Status,
		"payment_status": notification.Event.Transaction.Status,
		"request_id":     notification.Headers.RequestID,
	})

	if err := c.billingRepo.UpdatePaymentAttemptStatus(ctx, attempt.ID, normalized.PaymentAttemptStatus, map[string]any{
		"raw_last_status": string(rawLastStatusJSON),
	}); err != nil {
		return nil, fmt.Errorf("update payment attempt %s status: %w", attempt.ID, err)
	}

	if !normalized.IsTerminal {
		return &coreentity.DOKUWebhookResult{
			Provider:      "doku",
			InvoiceNumber: notification.Event.Order.InvoiceNumber,
			OrderStatus:   notification.Event.Order.Status,
			PaymentStatus: notification.Event.Transaction.Status,
		}, nil
	}

	if normalized.PaymentAttemptStatus != coreentity.PaymentAttemptStatusSucceeded {
		return &coreentity.DOKUWebhookResult{
			Provider:      "doku",
			InvoiceNumber: notification.Event.Order.InvoiceNumber,
			OrderStatus:   notification.Event.Order.Status,
			PaymentStatus: notification.Event.Transaction.Status,
		}, nil
	}

	invoice, err := c.billingRepo.GetInvoiceByNumber(ctx, notification.Event.Order.InvoiceNumber)
	if err != nil {
		return nil, fmt.Errorf("lookup invoice by number %s: %w", notification.Event.Order.InvoiceNumber, err)
	}

	if invoice.Status == coreentity.InvoiceStatusPaid {
		log.Ctx(ctx).Info().
			Str("invoice_number", invoice.InvoiceNumber).
			Msg("Invoice already paid, skipping subscription activation")
		return &coreentity.DOKUWebhookResult{
			Provider:      "doku",
			InvoiceNumber: notification.Event.Order.InvoiceNumber,
			OrderStatus:   notification.Event.Order.Status,
			PaymentStatus: notification.Event.Transaction.Status,
		}, nil
	}

	if err := c.billingRepo.UpdateInvoiceStatus(ctx, invoice.ID, coreentity.InvoiceStatusPaid); err != nil {
		return nil, fmt.Errorf("update invoice %s status: %w", invoice.ID, err)
	}

	targetState := invoice.TargetSubscriptionState
	if targetState == "" {
		targetState = coreentity.InternalTenantSubscriptionStatusPendingActivation
	}

	var subscriptionID string
	var accountID string

	existingSub, err := c.billingRepo.GetSubscription(ctx, invoice.TenantID)
	if err == nil && existingSub != nil {
		subscriptionID = existingSub.ID
		accountID = existingSub.InternalBillingAccountID
	}

	productID := ""
	pricingID := ""
	billingCycle := ""
	planName := ""
	productFeatures := []string{}
	if invoice.Metadata != nil {
		if pid, ok := invoice.Metadata["internal_product_id"].(string); ok {
			productID = pid
		}
		if pid, ok := invoice.Metadata["internal_product_pricing_id"].(string); ok {
			pricingID = pid
		}
		if bc, ok := invoice.Metadata["billing_cycle"].(string); ok {
			billingCycle = bc
		}
		if pn, ok := invoice.Metadata["internal_product_name"].(string); ok {
			planName = pn
		}
		if features, ok := invoice.Metadata["features"].([]interface{}); ok {
			for _, f := range features {
				if fs, ok := f.(string); ok {
					productFeatures = append(productFeatures, fs)
				}
			}
		}
	}

	now := time.Now().UTC()
	periodEnd := now.AddDate(0, 1, 0)
	if billingCycle == "annual" || billingCycle == "yearly" {
		periodEnd = now.AddDate(1, 0, 0)
	}
	periodStartStr := now.Format(time.RFC3339)
	periodEndStr := periodEnd.Format(time.RFC3339)

	planSnapshot := map[string]any{
		"name":          planName,
		"product_id":    productID,
		"pricing_id":    pricingID,
		"billing_cycle": billingCycle,
		"features":      productFeatures,
	}

	sub, err := c.billingRepo.UpsertSubscription(ctx, coreentity.InternalTenantSubscription{
		ID:                       subscriptionID,
		InternalBillingAccountID: accountID,
		TenantID:                 invoice.TenantID,
		ProductFamily:            "subscription",
		Status:                   targetState,
		InternalProductID:        productID,
		InternalProductPricingID: pricingID,
		PlanSnapshot:             planSnapshot,
		CurrentPeriodStartedAt:   &periodStartStr,
		CurrentPeriodEndedAt:     &periodEndStr,
		Metadata: map[string]any{
			"activated_from_invoice": invoice.ID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("upsert subscription for tenant %s: %w", invoice.TenantID, err)
	}

	_ = c.billingRepo.CreateSubscriptionChange(ctx, coreentity.InternalTenantSubscriptionChange{
		TenantID:                     invoice.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		ChangeType:                   coreentity.InternalTenantSubscriptionChangeTypeCheckout,
		FromStatus:                   coreentity.InternalTenantSubscriptionStatusFree,
		ToStatus:                     targetState,
		Trigger:                      coreentity.InternalTenantSubscriptionTriggerInvoicePaidFinal,
		SourceType:                   coreentity.InternalBillingSourceTypeDOKUWebhook,
		SourceReferenceID:            &invoice.ID,
		EffectiveAt:                  now.Format(time.RFC3339),
		Metadata: map[string]any{
			"invoice_number": invoice.InvoiceNumber,
			"provider":       "doku",
		},
	})

	_ = c.billingRepo.CreateEntitlementSnapshot(ctx, coreentity.InternalEntitlementSnapshot{
		TenantID:                     invoice.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		SubscriptionStatus:           targetState,
		Features:                     productFeatures,
		UsageLimits:                  map[string]int64{},
		SourceSnapshot:               planSnapshot,
		EffectiveAt:                  now.Format(time.RFC3339),
		Metadata: map[string]any{
			"source": "doku_webhook_activation",
		},
	})

	_ = c.billingRepo.CreateLedgerEntry(ctx, coreentity.InternalBillingLedgerEntry{
		TenantID:                     invoice.TenantID,
		InternalTenantSubscriptionID: &sub.ID,
		EntryType:                    coreentity.InternalBillingLedgerEntryTypePaymentSucceeded,
		SourceType:                   coreentity.InternalBillingSourceTypeDOKUWebhook,
		SourceReferenceID:            &invoice.ID,
		OccurredAt:                   now.Format(time.RFC3339),
		Metadata: map[string]any{
			"invoice_number": invoice.InvoiceNumber,
			"provider":       "doku",
		},
	})

	return &coreentity.DOKUWebhookResult{
		Provider:      "doku",
		InvoiceNumber: notification.Event.Order.InvoiceNumber,
		OrderStatus:   notification.Event.Order.Status,
		PaymentStatus: notification.Event.Transaction.Status,
	}, nil
}
