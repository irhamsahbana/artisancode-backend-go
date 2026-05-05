package core

import (
	"context"
	"sort"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) GetPlans(ctx context.Context) ([]coreentity.TenantBillingPlan, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_plans:GetPlans")
	defer span.End()

	if c.productRepo == nil || c.currencyRepo == nil {
		return []coreentity.TenantBillingPlan{}, nil
	}

	activeSubProductID := ""
	if c.billingRepo != nil {
		userCtx := common.GetUserContext(ctx)
		if userCtx.TenantID != "" {
			sub, err := c.billingRepo.GetActiveSubscription(ctx, userCtx.TenantID)
			if err == nil && sub != nil {
				activeSubProductID = sub.InternalProductID
			}
		}
	}

	products, _, err := c.productRepo.GetInternalProducts(ctx, coreentity.InternalProductListFilter{
		Status:   coreentity.InternalProductStatusActive,
		Page:     1,
		Paginate: 100,
	})
	if err != nil {
		return nil, err
	}

	plans := make([]coreentity.TenantBillingPlan, 0, len(products))
	for _, product := range products {
		pricings, _, err := c.productRepo.GetInternalProductPricings(ctx, coreentity.InternalProductPricingListFilter{
			InternalProductID: product.ID,
			Page:              1,
			Paginate:          100,
		})
		if err != nil {
			return nil, err
		}

		plan := coreentity.TenantBillingPlan{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Features:      stringSliceFromMetadata(product.Metadata, "features"),
			IsCurrentPlan: product.ID == activeSubProductID,
		}

		for _, pricing := range pricings {
			if pricing.Status != coreentity.InternalProductStatusActive {
				continue
			}

			cycle := billingCycleFromPricing(pricing)
			prices, _, err := c.productRepo.GetInternalProductPrices(ctx, coreentity.InternalProductPriceListFilter{
				InternalProductPricingID: pricing.ID,
				Page:                     1,
				Paginate:                 100,
			})
			if err != nil {
				return nil, err
			}

			for _, price := range prices {
				if !isPriceActiveNow(price) {
					continue
				}

				currency, err := c.currencyRepo.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{
					Code: price.CurrencyCode,
				})
				if err != nil || !currency.IsActive {
					continue
				}

				plan.Prices = append(plan.Prices, coreentity.TenantBillingPlanPrice{
					ID:                    price.ID,
					PricingID:             pricing.ID,
					BillingCycle:          cycle,
					Amount:                price.Amount.String(),
					Currency:              price.CurrencyCode,
					CurrencySymbol:        currency.Symbol,
					CurrencyDecimalPlaces: currency.DecimalPlaces,
					IsDefaultCurrency:     currency.IsDefault,
				})
			}
		}

		if len(plan.Prices) == 0 {
			continue
		}

		sort.SliceStable(plan.Prices, func(i, j int) bool {
			left := plan.Prices[i]
			right := plan.Prices[j]

			if left.IsDefaultCurrency != right.IsDefaultCurrency {
				return left.IsDefaultCurrency
			}
			if left.BillingCycle != right.BillingCycle {
				return left.BillingCycle == "monthly"
			}
			if left.Currency != right.Currency {
				return left.Currency < right.Currency
			}

			return left.Amount < right.Amount
		})

		firstPrice := plan.Prices[0]
		plan.PricingID = firstPrice.PricingID
		plan.BillingCycle = firstPrice.BillingCycle
		plan.Amount = firstPrice.Amount
		plan.Currency = firstPrice.Currency

		plans = append(plans, plan)
	}

	return plans, nil
}

func billingCycleFromPricing(pricing coreentity.InternalProductPricing) string {
	if value, ok := pricing.Metadata["billing_cycle"].(string); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}

	code := strings.ToLower(pricing.Code + " " + pricing.Name)
	if strings.Contains(code, "annual") || strings.Contains(code, "year") {
		return "annual"
	}

	return "monthly"
}

func stringSliceFromMetadata(metadata map[string]any, key string) []string {
	raw, ok := metadata[key]
	if !ok {
		return []string{}
	}

	values, ok := raw.([]any)
	if !ok {
		return []string{}
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		item, ok := value.(string)
		if ok && strings.TrimSpace(item) != "" {
			result = append(result, strings.TrimSpace(item))
		}
	}

	return result
}
