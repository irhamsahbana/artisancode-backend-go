package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"codebase-app/internal/entity/coreentity"
)

var ErrEntitlementPlanRequired = errors.New("entitlement plan required")

func CalculateEntitlementSnapshot(
	input coreentity.InternalEntitlementCalculationInput,
) (coreentity.InternalEntitlementSnapshot, error) {
	switch input.SubscriptionStatus {
	case coreentity.InternalTenantSubscriptionStatusActive,
		coreentity.InternalTenantSubscriptionStatusGracePeriod,
		coreentity.InternalTenantSubscriptionStatusCancelled:
		if input.Plan == nil {
			return coreentity.InternalEntitlementSnapshot{}, fmt.Errorf(
				"%w: %s",
				ErrEntitlementPlanRequired,
				input.SubscriptionStatus,
			)
		}

		features := featureSet(input.FreeTier.Features)
		limits := copyLimits(input.FreeTier.UsageLimits)
		mergeGrant(features, limits, *input.Plan)
		for _, addOn := range input.AddOns {
			mergeAddOn(features, limits, addOn)
		}

		return coreentity.InternalEntitlementSnapshot{
			SubscriptionStatus: input.SubscriptionStatus,
			Features:           sortedFeatures(features),
			UsageLimits:        limits,
		}, nil
	case coreentity.InternalTenantSubscriptionStatusSuspended:
		return coreentity.InternalEntitlementSnapshot{
			SubscriptionStatus: input.SubscriptionStatus,
			Features:           sortedFeatures(featureSet(input.SafeAccess.Features)),
			UsageLimits:        copyLimits(input.SafeAccess.UsageLimits),
		}, nil
	case coreentity.InternalTenantSubscriptionStatusFree,
		coreentity.InternalTenantSubscriptionStatusPendingActivation,
		coreentity.InternalTenantSubscriptionStatusExpired:
		return coreentity.InternalEntitlementSnapshot{
			SubscriptionStatus: input.SubscriptionStatus,
			Features:           sortedFeatures(featureSet(input.FreeTier.Features)),
			UsageLimits:        copyLimits(input.FreeTier.UsageLimits),
		}, nil
	default:
		return coreentity.InternalEntitlementSnapshot{}, fmt.Errorf(
			"%w: %s",
			ErrInvalidSubscriptionTransition,
			input.SubscriptionStatus,
		)
	}
}

func mergeGrant(
	features map[string]struct{},
	limits map[string]int64,
	grant coreentity.InternalEntitlementGrant,
) {
	for _, feature := range grant.Features {
		addFeature(features, feature)
	}
	for key, value := range grant.UsageLimits {
		limits[key] = value
	}
}

func mergeAddOn(
	features map[string]struct{},
	limits map[string]int64,
	addOn coreentity.InternalEntitlementAddOn,
) {
	switch addOn.Type {
	case coreentity.InternalEntitlementAddOnTypeFeatureUnlock:
		for _, feature := range addOn.Features {
			addFeature(features, feature)
		}
	case coreentity.InternalEntitlementAddOnTypeUsageLimitExtension:
		for key, value := range addOn.UsageLimits {
			limits[key] += value
		}
	default:
		for _, feature := range addOn.Features {
			addFeature(features, feature)
		}
		for key, value := range addOn.UsageLimits {
			limits[key] += value
		}
	}
}

func featureSet(features []string) map[string]struct{} {
	result := make(map[string]struct{}, len(features))
	for _, feature := range features {
		addFeature(result, feature)
	}
	return result
}

func addFeature(features map[string]struct{}, feature string) {
	feature = strings.TrimSpace(feature)
	if feature == "" {
		return
	}
	features[feature] = struct{}{}
}

func sortedFeatures(features map[string]struct{}) []string {
	result := make([]string, 0, len(features))
	for feature := range features {
		result = append(result, feature)
	}
	sort.Strings(result)
	return result
}

func copyLimits(limits map[string]int64) map[string]int64 {
	result := make(map[string]int64, len(limits))
	for key, value := range limits {
		result[key] = value
	}
	return result
}
