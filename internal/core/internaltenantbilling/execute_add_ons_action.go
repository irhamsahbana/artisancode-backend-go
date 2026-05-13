package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

type addOnMetadata struct {
	AddOnType       string           `json:"add_on_type"`
	Features        []string         `json:"features"`
	UsageLimits     map[string]int64 `json:"usage_limits"`
	CompatiblePlans []string         `json:"compatible_plans"`
}

func (c *internalTenantBillingCore) ExecuteAddOnsAction(
	ctx context.Context,
	input coreentity.TenantBillingAddOnsActionInput,
) (*coreentity.TenantBillingAddOnsActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:execute_add_ons_action:ExecuteAddOnsAction")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	input.UserCtx = userCtx

	switch input.Action {
	case coreentity.TenantBillingAddOnsActionAdd:
		return c.addAddOns(ctx, input)
	case coreentity.TenantBillingAddOnsActionRemove:
		return c.removeAddOns(ctx, input)
	default:
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedAddOnsAction))
	}
}

func (c *internalTenantBillingCore) addAddOns(
	ctx context.Context,
	input coreentity.TenantBillingAddOnsActionInput,
) (*coreentity.TenantBillingAddOnsActionResult, error) {
	if len(input.AddOnIDs) == 0 {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageAddOnsNotFound))
	}

	addOns, err := c.billingRepo.GetAddOnsByIDs(ctx, input.AddOnIDs)
	if err != nil {
		return nil, err
	}
	if len(addOns) != len(input.AddOnIDs) {
		return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageAddOnsNotFound))
	}

	sub, err := c.billingRepo.GetActiveSubscription(ctx, input.UserCtx.TenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageSubscriptionNotFound))
		}
		return nil, err
	}

	newSnapshots, err := buildAddOnSnapshots(addOns)
	if err != nil {
		return nil, err
	}

	baseProduct, err := c.productRepo.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: sub.InternalProductID})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", sub.ID).Msg("Failed to get base product for compatibility check")
		return nil, err
	}

	for _, addOn := range addOns {
		if addOn.RawMetadata == nil || len(*addOn.RawMetadata) == 0 {
			continue
		}
		var meta addOnMetadata
		if err := json.Unmarshal(*addOn.RawMetadata, &meta); err != nil {
			return nil, err
		}
		if len(meta.CompatiblePlans) == 0 {
			continue
		}
		found := false
		for _, planCode := range meta.CompatiblePlans {
			if planCode == baseProduct.Code {
				found = true
				break
			}
		}
		if !found {
			return nil, errmsg.NewCustomErrors(422,
				errmsg.WithMessage(errmsg.MessageAddOnsIncompatible),
			)
		}
	}

	existingIDs := make(map[string]bool, len(sub.AddOnSnapshots))
	for _, s := range sub.AddOnSnapshots {
		existingIDs[s.InternalProductID] = true
	}

	for _, s := range newSnapshots {
		if existingIDs[s.InternalProductID] {
			continue
		}
		sub.AddOnSnapshots = append(sub.AddOnSnapshots, s)
	}

	sub.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	sub, err = c.billingRepo.UpsertSubscription(ctx, *sub)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if err := c.billingRepo.CreateSubscriptionChange(ctx, coreentity.InternalTenantSubscriptionChange{
		TenantID:                     input.UserCtx.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		ChangeType:                   coreentity.InternalTenantSubscriptionChangeTypeUpgrade,
		FromStatus:                   sub.Status,
		ToStatus:                     sub.Status,
		Trigger:                      coreentity.InternalTenantSubscriptionTriggerAddOnAdded,
		SourceType:                   coreentity.InternalBillingSourceTypeSelfServeUpgrade,
		EffectiveAt:                  now,
		Metadata: map[string]any{
			"add_on_ids": input.AddOnIDs,
		},
	}); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("best-effort add-on subscription change record creation failed")
	}

	entitlement, err := buildEntitlementFromSubscription(*sub, nil)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("best-effort entitlement calculation failed")
	} else {
		if err := c.billingRepo.CreateEntitlementSnapshot(ctx, coreentity.InternalEntitlementSnapshot{
			TenantID:                     input.UserCtx.TenantID,
			InternalTenantSubscriptionID: sub.ID,
			SubscriptionStatus:           sub.Status,
			Features:                     entitlement.Features,
			UsageLimits:                  entitlement.UsageLimits,
			EffectiveAt:                  now,
			Metadata: map[string]any{
				"source":         "add_on_action",
				"add_on_ids":     input.AddOnIDs,
				"subscription_id": sub.ID,
			},
		}); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("best-effort add-on entitlement snapshot creation failed")
		}
	}

	_ = c.billingRepo.CreateLedgerEntry(ctx, coreentity.InternalBillingLedgerEntry{
		TenantID:                     input.UserCtx.TenantID,
		InternalTenantSubscriptionID: &sub.ID,
		EntryType:                    coreentity.InternalBillingLedgerEntryTypeSubscriptionChange,
		SourceType:                   coreentity.InternalBillingSourceTypeSelfServeUpgrade,
		OccurredAt:                   now,
		Metadata: map[string]any{
			"add_on_ids": input.AddOnIDs,
			"action":     "add",
		},
	})

	resultIDs := make([]string, 0, len(sub.AddOnSnapshots))
	for _, s := range sub.AddOnSnapshots {
		resultIDs = append(resultIDs, s.InternalProductID)
	}

	return &coreentity.TenantBillingAddOnsActionResult{
		SubscriptionID: sub.ID,
		AddOnIDs:       resultIDs,
	}, nil
}

func (c *internalTenantBillingCore) removeAddOns(
	ctx context.Context,
	input coreentity.TenantBillingAddOnsActionInput,
) (*coreentity.TenantBillingAddOnsActionResult, error) {
	if len(input.AddOnIDs) == 0 {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageAddOnsNotFound))
	}

	sub, err := c.billingRepo.GetActiveSubscription(ctx, input.UserCtx.TenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageSubscriptionNotFound))
		}
		return nil, err
	}

	removeSet := make(map[string]bool, len(input.AddOnIDs))
	for _, id := range input.AddOnIDs {
		removeSet[id] = true
	}

	filtered := make([]coreentity.InternalTenantSubscriptionAddOnSnapshot, 0, len(sub.AddOnSnapshots))
	for _, s := range sub.AddOnSnapshots {
		if removeSet[s.InternalProductID] {
			continue
		}
		filtered = append(filtered, s)
	}
	sub.AddOnSnapshots = filtered

	sub.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	sub, err = c.billingRepo.UpsertSubscription(ctx, *sub)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if err := c.billingRepo.CreateSubscriptionChange(ctx, coreentity.InternalTenantSubscriptionChange{
		TenantID:                     input.UserCtx.TenantID,
		InternalTenantSubscriptionID: sub.ID,
		ChangeType:                   coreentity.InternalTenantSubscriptionChangeTypeDowngrade,
		FromStatus:                   sub.Status,
		ToStatus:                     sub.Status,
		Trigger:                      coreentity.InternalTenantSubscriptionTriggerAddOnRemoved,
		SourceType:                   coreentity.InternalBillingSourceTypeSelfServeUpgrade,
		EffectiveAt:                  now,
		Metadata: map[string]any{
			"add_on_ids": input.AddOnIDs,
		},
	}); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("best-effort add-on removal subscription change record creation failed")
	}

	entitlement, err := buildEntitlementFromSubscription(*sub, nil)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("best-effort entitlement calculation failed")
	} else {
		if err := c.billingRepo.CreateEntitlementSnapshot(ctx, coreentity.InternalEntitlementSnapshot{
			TenantID:                     input.UserCtx.TenantID,
			InternalTenantSubscriptionID: sub.ID,
			SubscriptionStatus:           sub.Status,
			Features:                     entitlement.Features,
			UsageLimits:                  entitlement.UsageLimits,
			EffectiveAt:                  now,
			Metadata: map[string]any{
				"source":         "remove_add_on_action",
				"add_on_ids":     input.AddOnIDs,
				"subscription_id": sub.ID,
			},
		}); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("best-effort add-on removal entitlement snapshot creation failed")
		}
	}

	_ = c.billingRepo.CreateLedgerEntry(ctx, coreentity.InternalBillingLedgerEntry{
		TenantID:                     input.UserCtx.TenantID,
		InternalTenantSubscriptionID: &sub.ID,
		EntryType:                    coreentity.InternalBillingLedgerEntryTypeSubscriptionChange,
		SourceType:                   coreentity.InternalBillingSourceTypeSelfServeUpgrade,
		OccurredAt:                   now,
		Metadata: map[string]any{
			"add_on_ids": input.AddOnIDs,
			"action":     "remove",
		},
	})

	resultIDs := make([]string, 0, len(sub.AddOnSnapshots))
	for _, s := range sub.AddOnSnapshots {
		resultIDs = append(resultIDs, s.InternalProductID)
	}

	return &coreentity.TenantBillingAddOnsActionResult{
		SubscriptionID: sub.ID,
		AddOnIDs:       resultIDs,
	}, nil
}

func buildAddOnSnapshots(
	addOns []coreentity.TenantBillingAddOn,
) ([]coreentity.InternalTenantSubscriptionAddOnSnapshot, error) {
	snapshots := make([]coreentity.InternalTenantSubscriptionAddOnSnapshot, 0, len(addOns))

	for _, addOn := range addOns {
		if addOn.RawMetadata == nil || len(*addOn.RawMetadata) == 0 {
			snapshots = append(snapshots, coreentity.InternalTenantSubscriptionAddOnSnapshot{
				InternalProductID:        addOn.ID,
				InternalProductPricingID: addOn.PricingID,
			})
			continue
		}

		var meta addOnMetadata
		if err := json.Unmarshal(*addOn.RawMetadata, &meta); err != nil {
			return nil, err
		}

		snapshots = append(snapshots, coreentity.InternalTenantSubscriptionAddOnSnapshot{
			InternalProductID:        addOn.ID,
			InternalProductPricingID: addOn.PricingID,
			Type:                     meta.AddOnType,
			Features:                 meta.Features,
			UsageLimits:              meta.UsageLimits,
		})
	}

	return snapshots, nil
}

func buildEntitlementFromSubscription(
	sub coreentity.InternalTenantSubscription,
	planGrant *coreentity.InternalEntitlementGrant,
) (coreentity.InternalEntitlementSnapshot, error) {
	addOns := make([]coreentity.InternalEntitlementAddOn, 0, len(sub.AddOnSnapshots))
	for _, snap := range sub.AddOnSnapshots {
		addOns = append(addOns, coreentity.InternalEntitlementAddOn{
			Type:        snap.Type,
			Features:    snap.Features,
			UsageLimits: snap.UsageLimits,
		})
	}

	input := coreentity.InternalEntitlementCalculationInput{
		SubscriptionStatus: sub.Status,
		FreeTier: coreentity.InternalEntitlementGrant{
			Features:    []string{"attendance_basic"},
			UsageLimits: map[string]int64{"employees": 10, "branches": 1},
		},
		SafeAccess: coreentity.InternalEntitlementGrant{
			Features:    []string{"account_recovery"},
			UsageLimits: map[string]int64{"employees": 0, "branches": 0},
		},
		Plan:   planGrant,
		AddOns: addOns,
	}

	return CalculateEntitlementSnapshot(input)
}
