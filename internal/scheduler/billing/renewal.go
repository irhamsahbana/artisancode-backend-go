package billing

import (
	"context"
	"fmt"
	"time"

	core "codebase-app/internal/core/internaltenantbilling"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	repository "codebase-app/internal/ports/secondary/db"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	repo       repository.InternalTenantBillingRepository
	batchLimit int
	graceDays  int
}

func New(repo repository.InternalTenantBillingRepository, batchLimit, graceDays int) *Handler {
	if batchLimit <= 0 {
		batchLimit = 50
	}
	if graceDays <= 0 {
		graceDays = 14
	}
	return &Handler{repo: repo, batchLimit: batchLimit, graceDays: graceDays}
}

func (h *Handler) ProcessRenewals(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "internal:scheduler:billing:process_renewals")
	defer span.End()

	subs, err := h.repo.GetSubscriptionsPastPeriodEnd(ctx, h.batchLimit)
	if err != nil {
		return fmt.Errorf("get subscriptions past period end: %w", err)
	}

	log.Ctx(ctx).Info().Int("count", len(subs)).Msg("Processing subscription renewals")

	for _, sub := range subs {
		h.processSingleRenewal(ctx, sub)
	}

	return nil
}

func (h *Handler) processSingleRenewal(ctx context.Context, sub coreentity.InternalTenantSubscription) {
	if !core.CanTransitionSubscription(sub.Status, coreentity.InternalTenantSubscriptionTriggerRenewalUnpaid) {
		log.Ctx(ctx).Warn().Str("subscription_id", sub.ID).Str("status", sub.Status).Msg("Cannot transition subscription for renewal")
		return
	}

	nextStatus, err := core.NextSubscriptionStatus(sub.Status, coreentity.InternalTenantSubscriptionTriggerRenewalUnpaid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to get next subscription status")
		return
	}

	pricing, err := h.repo.GetPricingInfo(ctx, sub.InternalProductPricingID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Str("pricing_id", sub.InternalProductPricingID).Msg("Failed to get pricing info")
		return
	}

	now := time.Now().UTC()
	periodStart := now.Format(time.RFC3339)
	periodEnd := now.AddDate(0, 1, 0).Format(time.RFC3339)
	if pricing.BillingCycle == "annual" {
		periodEnd = now.AddDate(1, 0, 0).Format(time.RFC3339)
	}
	graceEnd := now.AddDate(0, 0, h.graceDays)
	graceEndStr := graceEnd.Format(time.RFC3339)

	invNumber := fmt.Sprintf("TINV-%s", uuid.NewString()[:8])

	invoice, err := h.repo.CreateInternalInvoice(ctx, coreentity.InternalTenantInvoice{
		InvoiceNumber:                 invNumber,
		TenantID:                      sub.TenantID,
		InternalTenantSubscriptionID:  &sub.ID,
		Status:                        coreentity.InternalTenantInvoiceStatusPending,
		CurrencyCode:                  pricing.CurrencyCode,
		Amount:                        pricing.Amount,
		AmountPaid:                    "0",
		AmountOutstanding:             pricing.Amount,
		SourceType:                    coreentity.InternalBillingSourceTypeRenewalScheduler,
		DueAt:                         &graceEndStr,
		Metadata: map[string]any{
			"renewal":         true,
			"subscription_id": sub.ID,
			"billing_cycle":   pricing.BillingCycle,
		},
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to create renewal invoice")
		return
	}

	sub.Status = nextStatus
	sub.CurrentPeriodStartedAt = &periodStart
	sub.CurrentPeriodEndedAt = &periodEnd
	sub.GraceEndedAt = &graceEndStr

	_, err = h.repo.UpsertSubscription(ctx, sub)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to update subscription for renewal")
		return
	}

	effectiveAt := now.Format(time.RFC3339)
	_ = h.repo.CreateSubscriptionChange(ctx, coreentity.InternalTenantSubscriptionChange{
		TenantID:                     sub.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		ChangeType:                   coreentity.InternalTenantSubscriptionChangeTypeUpgrade,
		FromStatus:                   coreentity.InternalTenantSubscriptionStatusActive,
		ToStatus:                     nextStatus,
		Trigger:                      coreentity.InternalTenantSubscriptionTriggerRenewalUnpaid,
		SourceType:                   coreentity.InternalBillingSourceTypeRenewalScheduler,
		EffectiveAt:                  effectiveAt,
		Metadata: map[string]any{
			"invoice_id":    invoice.ID,
			"billing_cycle": pricing.BillingCycle,
		},
	})

	_ = h.repo.CreateLedgerEntry(ctx, coreentity.InternalBillingLedgerEntry{
		TenantID:                     sub.TenantID,
		InternalTenantSubscriptionID: &sub.ID,
		EntryType:                    coreentity.InternalBillingLedgerEntryTypeSubscriptionChange,
		SourceType:                   coreentity.InternalBillingSourceTypeRenewalScheduler,
		OccurredAt:                   effectiveAt,
		Metadata: map[string]any{
			"invoice_id": invoice.ID,
			"renewal":    true,
		},
	})

	log.Ctx(ctx).Info().
		Str("subscription_id", sub.ID).
		Str("tenant_id", sub.TenantID).
		Str("invoice_id", invoice.ID).
		Str("from_status", coreentity.InternalTenantSubscriptionStatusActive).
		Str("to_status", nextStatus).
		Msg("Subscription renewed with unpaid invoice")
}
