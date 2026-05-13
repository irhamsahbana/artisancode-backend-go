package billing

import (
	"context"
	"fmt"
	"time"

	core "codebase-app/internal/core/internaltenantbilling"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (h *Handler) ProcessDunningEscalation(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "internal:scheduler:billing:process_dunning_escalation")
	defer span.End()

	subs, err := h.repo.GetSubscriptionsInGracePastDue(ctx, h.batchLimit)
	if err != nil {
		return fmt.Errorf("get subscriptions in grace past due: %w", err)
	}

	log.Ctx(ctx).Info().Int("count", len(subs)).Msg("Processing dunning escalations")

	for _, sub := range subs {
		h.processSingleDunning(ctx, sub)
	}

	return nil
}

func (h *Handler) processSingleDunning(ctx context.Context, sub coreentity.InternalTenantSubscription) {
	if !core.CanTransitionSubscription(sub.Status, coreentity.InternalTenantSubscriptionTriggerGraceEnded) {
		log.Ctx(ctx).Warn().Str("subscription_id", sub.ID).Str("status", sub.Status).Msg("Cannot transition subscription for dunning escalation")
		return
	}

	nextStatus, err := core.NextSubscriptionStatus(sub.Status, coreentity.InternalTenantSubscriptionTriggerGraceEnded)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to get next subscription status")
		return
	}

	now := time.Now().UTC()
	sub.Status = nextStatus
	sub.UpdatedAt = now.Format(time.RFC3339)

	_, err = h.repo.UpsertSubscription(ctx, sub)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to update subscription for dunning escalation")
		return
	}

	effectiveAt := now.Format(time.RFC3339)
	_ = h.repo.CreateSubscriptionChange(ctx, coreentity.InternalTenantSubscriptionChange{
		TenantID:                     sub.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		ChangeType:                   coreentity.InternalTenantSubscriptionChangeTypeDowngrade,
		FromStatus:                   coreentity.InternalTenantSubscriptionStatusGracePeriod,
		ToStatus:                     nextStatus,
		Trigger:                      coreentity.InternalTenantSubscriptionTriggerGraceEnded,
		SourceType:                   coreentity.InternalBillingSourceTypeRenewalScheduler,
		EffectiveAt:                  effectiveAt,
	})

	_ = h.repo.CreateLedgerEntry(ctx, coreentity.InternalBillingLedgerEntry{
		TenantID:                     sub.TenantID,
		InternalTenantSubscriptionID: &sub.ID,
		EntryType:                    coreentity.InternalBillingLedgerEntryTypeSubscriptionChange,
		SourceType:                   coreentity.InternalBillingSourceTypeRenewalScheduler,
		OccurredAt:                   effectiveAt,
		Metadata: map[string]any{
			"dunning_escalation": true,
			"grace_ended":        true,
		},
	})

	log.Ctx(ctx).Info().
		Str("subscription_id", sub.ID).
		Str("tenant_id", sub.TenantID).
		Str("from_status", coreentity.InternalTenantSubscriptionStatusGracePeriod).
		Str("to_status", nextStatus).
		Msg("Subscription suspended due to dunning escalation")
}
