package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func TenantBillingPlanFromCoreToRest(item coreentity.TenantBillingPlan) restentity.TenantBillingPlan {
	prices := make([]restentity.TenantBillingPlanPrice, 0, len(item.Prices))
	for _, price := range item.Prices {
		prices = append(prices, restentity.TenantBillingPlanPrice{
			ID:                    price.ID,
			PricingID:             price.PricingID,
			BillingCycle:          price.BillingCycle,
			Amount:                price.Amount,
			Currency:              price.Currency,
			CurrencySymbol:        price.CurrencySymbol,
			CurrencyDecimalPlaces: price.CurrencyDecimalPlaces,
			IsDefaultCurrency:     price.IsDefaultCurrency,
		})
	}

	addOns := make([]restentity.TenantBillingAddOn, 0, len(item.AddOns))
	for _, addOn := range item.AddOns {
		addOns = append(addOns, tenantBillingAddOnFromCoreToRest(addOn))
	}

	return restentity.TenantBillingPlan{
		ID:           item.ID,
		Name:         item.Name,
		Description:  item.Description,
		PricingID:    item.PricingID,
		BillingCycle: item.BillingCycle,
		Amount:       item.Amount,
		Currency:     item.Currency,
		Features:     item.Features,
		Prices:       prices,
		AddOns:       addOns,
	}
}

func TenantBillingSubscriptionFromCoreToRest(
	item coreentity.TenantBillingSubscriptionView,
) restentity.TenantBillingSubscription {
	addOns := make([]restentity.TenantBillingAddOn, 0, len(item.AddOns))
	for _, addOn := range item.AddOns {
		addOns = append(addOns, tenantBillingAddOnFromCoreToRest(addOn))
	}

	return restentity.TenantBillingSubscription{
		ID:                 item.ID,
		Status:             item.Status,
		PlanName:           item.PlanName,
		BillingCycle:       item.BillingCycle,
		CurrentPeriodStart: item.CurrentPeriodStart,
		CurrentPeriodEnd:   item.CurrentPeriodEnd,
		NextRenewalAt:      item.NextRenewalAt,
		AddOns:             addOns,
	}
}

func TenantBillingEntitlementsFromCoreToRest(
	item coreentity.InternalEntitlementSnapshot,
) restentity.TenantBillingEntitlements {
	return restentity.TenantBillingEntitlements{
		SubscriptionState: item.SubscriptionStatus,
		Features:          item.Features,
		UsageLimits:       item.UsageLimits,
	}
}

func TenantBillingInvoiceFromCoreToRest(item coreentity.TenantBillingInvoice) restentity.TenantBillingInvoice {
	return restentity.TenantBillingInvoice{
		ID:               item.ID,
		InvoiceNumber:    item.InvoiceNumber,
		Status:           item.Status,
		Amount:           item.Amount,
		Currency:         item.Currency,
		PaymentAttemptID: item.PaymentAttemptID,
		PaymentURL:       item.PaymentURL,
		DueAt:            item.DueAt,
		PaidAt:           item.PaidAt,
		CreatedAt:        item.CreatedAt,
	}
}

func TenantBillingCheckoutFromRest(
	ctx context.Context,
	req restentity.CreateTenantBillingCheckoutReq,
	requestID string,
) coreentity.TenantBillingCheckoutInput {
	return coreentity.TenantBillingCheckoutInput{
		UserCtx:               common.GetUserContext(ctx),
		PriceID:               req.PriceID,
		AddOnIDs:              req.AddOnIDs,
		ReplaceActiveCheckout: req.ReplaceActiveCheckout,
		CallbackURL:           req.CallbackURL,
		RequestID:             requestID,
	}
}

func TenantBillingCheckoutFromCoreToRest(
	item coreentity.TenantBillingCheckoutResult,
) restentity.CreateTenantBillingCheckoutResp {
	return restentity.CreateTenantBillingCheckoutResp{
		InvoiceID:               item.InvoiceID,
		PaymentAttemptID:        item.PaymentAttemptID,
		PaymentURL:              item.PaymentURL,
		TargetSubscriptionState: item.TargetSubscriptionState,
		CheckoutReused:          item.CheckoutReused,
	}
}

func tenantBillingAddOnFromCoreToRest(item coreentity.TenantBillingAddOn) restentity.TenantBillingAddOn {
	return restentity.TenantBillingAddOn{
		ID:           item.ID,
		Name:         item.Name,
		Description:  item.Description,
		Amount:       item.Amount,
		Currency:     item.Currency,
		BillingCycle: item.BillingCycle,
	}
}
