package core

import (
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

func TestCalculateEntitlementSnapshot(t *testing.T) {
	freeTier := coreentity.InternalEntitlementGrant{
		Features: []string{"attendance_basic"},
		UsageLimits: map[string]int64{
			"employees": 10,
			"branches":  1,
		},
	}
	safeAccess := coreentity.InternalEntitlementGrant{
		Features: []string{"account_recovery"},
		UsageLimits: map[string]int64{
			"employees": 0,
			"branches":  0,
		},
	}
	plan := coreentity.InternalEntitlementGrant{
		Features: []string{"attendance_basic", "attendance_advanced"},
		UsageLimits: map[string]int64{
			"employees": 100,
			"branches":  5,
		},
	}

	tests := []struct {
		name     string
		input    coreentity.InternalEntitlementCalculationInput
		features []string
		limits   map[string]int64
	}{
		{
			name: "free uses explicit free tier",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusFree,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
			},
			features: []string{"attendance_basic"},
			limits: map[string]int64{
				"employees": 10,
				"branches":  1,
			},
		},
		{
			name: "pending activation does not grant paid features",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusPendingActivation,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
			},
			features: []string{"attendance_basic"},
			limits: map[string]int64{
				"employees": 10,
				"branches":  1,
			},
		},
		{
			name: "active applies plan and add-ons",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
				AddOns: []coreentity.InternalEntitlementAddOn{
					{
						Type:     coreentity.InternalEntitlementAddOnTypeFeatureUnlock,
						Features: []string{"payroll_export"},
					},
					{
						Type: coreentity.InternalEntitlementAddOnTypeUsageLimitExtension,
						UsageLimits: map[string]int64{
							"employees": 50,
						},
					},
				},
			},
			features: []string{"attendance_advanced", "attendance_basic", "payroll_export"},
			limits: map[string]int64{
				"employees": 150,
				"branches":  5,
			},
		},
		{
			name: "grace keeps paid entitlements",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusGracePeriod,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
			},
			features: []string{"attendance_advanced", "attendance_basic"},
			limits: map[string]int64{
				"employees": 100,
				"branches":  5,
			},
		},
		{
			name: "cancelled keeps paid entitlements until period end",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusCancelled,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
			},
			features: []string{"attendance_advanced", "attendance_basic"},
			limits: map[string]int64{
				"employees": 100,
				"branches":  5,
			},
		},
		{
			name: "suspended removes paid features",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusSuspended,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
			},
			features: []string{"account_recovery"},
			limits: map[string]int64{
				"employees": 0,
				"branches":  0,
			},
		},
		{
			name: "expired falls back to free tier",
			input: coreentity.InternalEntitlementCalculationInput{
				SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusExpired,
				FreeTier:           freeTier,
				SafeAccess:         safeAccess,
				Plan:               &plan,
			},
			features: []string{"attendance_basic"},
			limits: map[string]int64{
				"employees": 10,
				"branches":  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateEntitlementSnapshot(tt.input)

			require.NoError(t, err)
			require.Equal(t, tt.input.SubscriptionStatus, got.SubscriptionStatus)
			require.Equal(t, tt.features, got.Features)
			require.Equal(t, tt.limits, got.UsageLimits)
		})
	}
}

func TestCalculateEntitlementSnapshotRequiresPlanForPaidStates(t *testing.T) {
	got, err := CalculateEntitlementSnapshot(coreentity.InternalEntitlementCalculationInput{
		SubscriptionStatus: coreentity.InternalTenantSubscriptionStatusActive,
	})

	require.Empty(t, got)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrEntitlementPlanRequired))
}

func TestCalculateEntitlementSnapshotRejectsUnknownStatus(t *testing.T) {
	got, err := CalculateEntitlementSnapshot(coreentity.InternalEntitlementCalculationInput{
		SubscriptionStatus: "unknown",
	})

	require.Empty(t, got)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidSubscriptionTransition))
}
