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

func TenantBillingInvoiceDetailFromCoreToRest(
	item coreentity.TenantBillingInvoiceDetail,
) restentity.TenantBillingInvoiceDetail {
	attempts := make([]restentity.TenantBillingPaymentAttempt, 0, len(item.PaymentAttempts))
	for _, a := range item.PaymentAttempts {
		attempts = append(attempts, TenantBillingPaymentAttemptFromCoreToRest(a))
	}

	return restentity.TenantBillingInvoiceDetail{
		ID:               item.ID,
		InvoiceNumber:    item.InvoiceNumber,
		Status:           item.Status,
		Amount:           item.Amount,
		Currency:         item.CurrencyCode,
		AmountPaid:       item.AmountPaid,
		AmountOutstanding: item.AmountOutstanding,
		SourceType:       item.SourceType,
		DueAt:            item.DueAt,
		PaidAt:           item.PaidAt,
		ExpiredAt:        item.ExpiredAt,
		CreatedAt:        item.CreatedAt,
		PaymentAttempts:  attempts,
	}
}

func TenantBillingPaymentAttemptFromCoreToRest(
	item coreentity.TenantBillingPaymentAttemptView,
) restentity.TenantBillingPaymentAttempt {
	return restentity.TenantBillingPaymentAttempt{
		ID:                 item.ID,
		Provider:           item.Provider,
		PaymentMethodType:  item.PaymentMethodType,
		PaymentChannelCode: item.PaymentChannelCode,
		ProviderReference:  item.ProviderReference,
		PaymentURL:         item.PaymentURL,
		Status:             item.Status,
		RequestedAmount:    item.RequestedAmount,
		PaidAmount:         item.PaidAmount,
		ExpiredAt:          item.ExpiredAt,
		PaidAt:             item.PaidAt,
		FailedAt:           item.FailedAt,
		CreatedAt:          item.CreatedAt,
	}
}

func TenantBillingInvoiceActionFromRest(
	req restentity.ExecuteTenantBillingInvoiceActionReq,
) coreentity.TenantBillingInvoiceActionInput {
	return coreentity.TenantBillingInvoiceActionInput{
		ID:     req.ID,
		Action: req.Action,
		Reason: req.Reason,
	}
}

func TenantBillingInvoiceActionResultFromCoreToRest(
	item coreentity.TenantBillingInvoiceActionResult,
) restentity.ExecuteTenantBillingInvoiceActionResp {
	return restentity.ExecuteTenantBillingInvoiceActionResp{
		InvoiceID: item.InvoiceID,
		Status:    item.Status,
	}
}

func TenantBillingPaymentAttemptActionFromRest(
	req restentity.ExecuteTenantBillingPaymentAttemptActionReq,
) coreentity.TenantBillingPaymentAttemptActionInput {
	return coreentity.TenantBillingPaymentAttemptActionInput{
		ID:     req.ID,
		Action: req.Action,
	}
}

func TenantBillingPaymentAttemptActionResultFromCoreToRest(
	item coreentity.TenantBillingPaymentAttemptActionResult,
) restentity.ExecuteTenantBillingPaymentAttemptActionResp {
	return restentity.ExecuteTenantBillingPaymentAttemptActionResp{
		PaymentAttemptID: item.PaymentAttemptID,
		Status:           item.Status,
		PaymentURL:       item.PaymentURL,
	}
}

func TenantBillingSubscriptionActionFromRest(
	req restentity.ExecuteTenantBillingSubscriptionActionReq,
) coreentity.TenantBillingSubscriptionActionInput {
	return coreentity.TenantBillingSubscriptionActionInput{
		Action: req.Action,
	}
}

func TenantBillingSubscriptionActionResultFromCoreToRest(
	item coreentity.TenantBillingSubscriptionActionResult,
) restentity.ExecuteTenantBillingSubscriptionActionResp {
	return restentity.ExecuteTenantBillingSubscriptionActionResp{
		SubscriptionID: item.SubscriptionID,
		Status:         item.Status,
	}
}

func TenantBillingAddOnsActionFromRest(
	req restentity.ExecuteTenantBillingAddOnsActionReq,
) coreentity.TenantBillingAddOnsActionInput {
	return coreentity.TenantBillingAddOnsActionInput{
		Action:   req.Action,
		AddOnIDs: req.AddOnIDs,
	}
}

func TenantBillingAddOnsActionResultFromCoreToRest(
	item coreentity.TenantBillingAddOnsActionResult,
) restentity.ExecuteTenantBillingAddOnsActionResp {
	return restentity.ExecuteTenantBillingAddOnsActionResp{
		SubscriptionID: item.SubscriptionID,
		AddOnIDs:       item.AddOnIDs,
	}
}
