package core

import (
	"context"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

func (c *internalTenantBillingCore) CreateCheckout(
	ctx context.Context,
	input coreentity.TenantBillingCheckoutInput,
) (*coreentity.TenantBillingCheckoutResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:create_checkout:CreateCheckout")
	defer span.End()

	if !canManageBilling(input.UserCtx) {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageInvalidCredentials)
	}
	if c.billingRepo == nil || c.productRepo == nil || c.currencyRepo == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageDokuClientIsNotConfigured)
	}
	if c.doku == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageDokuClientIsNotConfigured)
	}

	price, err := c.productRepo.GetInternalProductPrice(ctx, coreentity.InternalProductPriceFilter{
		ID: input.PriceID,
	})
	if err != nil {
		return nil, err
	}
	if !isPriceActiveNow(*price) {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePriceIsNotActive)
	}

	isCurrencyActive, err := c.currencyRepo.IsCurrencyActive(ctx, price.CurrencyCode)
	if err != nil {
		return nil, err
	}
	if !isCurrencyActive {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyIsNotActive)
	}

	currency, err := c.currencyRepo.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{
		Code: price.CurrencyCode,
	})
	if err != nil {
		return nil, err
	}
	if currency.DecimalPlaces == 0 && !price.Amount.Equal(price.Amount.Truncate(0)) {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyAmountPrecisionIsInvalid)
	}

	pricing, err := c.productRepo.GetInternalProductPricing(ctx, coreentity.InternalProductPricingFilter{
		ID: price.InternalProductPricingID,
	})
	if err != nil {
		return nil, err
	}
	if pricing.Status != coreentity.InternalProductStatusActive {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePriceIsNotActive)
	}
	product, err := c.productRepo.GetInternalProduct(ctx, coreentity.InternalProductFilter{
		ID: pricing.InternalProductID,
	})
	if err != nil {
		return nil, err
	}
	if product.Status != coreentity.InternalProductStatusActive {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePriceIsNotActive)
	}

	var (
		invoice *coreentity.InternalTenantInvoice
		attempt *coreentity.InternalTenantPaymentAttempt
	)

	run := func(txCtx context.Context) error {
		account, err := c.billingRepo.EnsureBillingAccount(txCtx, input.UserCtx.TenantID)
		if err != nil {
			return err
		}

		invoice, err = c.billingRepo.CreateInvoice(txCtx, coreentity.InternalTenantInvoice{
			TenantID:                 input.UserCtx.TenantID,
			InternalBillingAccountID: account.ID,
			Status:                   coreentity.InvoiceStatusOpen,
			CurrencyCode:             price.CurrencyCode,
			Amount:                   price.Amount.String(),
			AmountPaid:               decimal.Zero.String(),
			AmountOutstanding:        price.Amount.String(),
			SourceType:               coreentity.InternalBillingSourceTypeSelfServeCheckout,
			TargetSubscriptionState:  coreentity.InternalTenantSubscriptionStatusPendingActivation,
			Metadata: map[string]any{
				"internal_product_id":         product.ID,
				"internal_product_name":       product.Name,
				"internal_product_price_id":   price.ID,
				"internal_product_pricing_id": pricing.ID,
				"billing_cycle":               billingCycleFromPricing(*pricing),
				"add_on_ids":                  input.AddOnIDs,
				"request_id":                  input.RequestID,
			},
		})
		if err != nil {
			return err
		}

		attempt, err = c.billingRepo.CreatePaymentAttempt(txCtx, coreentity.InternalTenantPaymentAttempt{
			TenantID:                input.UserCtx.TenantID,
			InternalTenantInvoiceID: invoice.ID,
			Provider:                coreentity.PaymentProviderDOKU,
			PaymentMethodType:       "gateway",
			ProviderReference:       invoice.InvoiceNumber,
			Status:                  coreentity.PaymentAttemptStatusInitiated,
			RequestedAmount:         invoice.Amount,
			PaidAmount:              decimal.Zero.String(),
			Metadata:                map[string]any{},
		})
		return err
	}

	if c.tx != nil {
		if err := c.tx.WithinTransaction(ctx, run); err != nil {
			return nil, err
		}
	} else if err := run(ctx); err != nil {
		return nil, err
	}

	amount := decimal.RequireFromString(invoice.Amount)
	multiplier := decimal.NewFromInt(1)
	if currency.DecimalPlaces > 0 {
		multiplier = decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(currency.DecimalPlaces)))
	}
	dokuAmount := amount.Mul(multiplier).Round(0).IntPart()

	dokuResp, err := c.doku.CreatePayment(ctx, restentity.DokuCreatePaymentRequest{
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        dokuAmount,
		Currency:      price.CurrencyCode,
		CustomerName:  checkoutCustomerName(input.UserCtx),
		CustomerEmail: checkoutCustomerEmail(input.UserCtx),
		LineItems: []restentity.DokuLineItem{
			{
				Name:     buildCheckoutLineItemName(product.Name, pricing.Name),
				Price:    dokuAmount,
				Quantity: 1,
			},
		},
		CallbackURL: input.CallbackURL,
	})
	if err != nil {
		return nil, err
	}

	attempt.ProviderRequestID = dokuResp.RequestID
	attempt.ProviderPaymentURL = dokuResp.PaymentURL
	attempt.ProviderPayloadSnapshot = map[string]any{"doku_response": dokuResp}
	attempt.Status = coreentity.PaymentAttemptStatusPending
	attempt, err = c.billingRepo.UpdatePaymentAttemptGateway(ctx, *attempt)
	if err != nil {
		return nil, err
	}

	return &coreentity.TenantBillingCheckoutResult{
		InvoiceID:               invoice.ID,
		PaymentAttemptID:        attempt.ID,
		PaymentURL:              attempt.ProviderPaymentURL,
		TargetSubscriptionState: coreentity.InternalTenantSubscriptionStatusPendingActivation,
		CheckoutReused:          false,
	}, nil
}

func canManageBilling(userCtx common.UserContext) bool {
	return userCtx.HasRole("owner") || userCtx.HasRole("admin")
}

func checkoutCustomerName(userCtx common.UserContext) string {
	if strings.TrimSpace(userCtx.UserName) != "" {
		return strings.TrimSpace(userCtx.UserName)
	}
	if strings.TrimSpace(userCtx.TenantName) != "" {
		return strings.TrimSpace(userCtx.TenantName)
	}
	return "Billing Admin"
}

func checkoutCustomerEmail(userCtx common.UserContext) string {
	if strings.TrimSpace(userCtx.TenantID) == "" {
		return "billing@example.com"
	}
	return "billing+" + strings.TrimSpace(userCtx.TenantID) + "@example.com"
}

func buildCheckoutLineItemName(productName, pricingName string) string {
	productName = strings.TrimSpace(productName)
	pricingName = strings.TrimSpace(pricingName)

	switch {
	case productName == "" && pricingName == "":
		return "Subscription"
	case productName == "":
		return pricingName
	case pricingName == "":
		return productName
	case strings.Contains(strings.ToLower(pricingName), strings.ToLower(productName)):
		return pricingName
	default:
		return productName + " - " + pricingName
	}
}

func isPriceActiveNow(price coreentity.InternalProductPrice) bool {
	now := time.Now().UTC()
	startedAt, err := time.Parse(time.RFC3339, price.StartedAt)
	if err != nil || startedAt.After(now) {
		return false
	}
	if price.EndedAt == nil {
		return true
	}
	endedAt, err := time.Parse(time.RFC3339, *price.EndedAt)
	if err != nil {
		return false
	}

	return endedAt.After(now)
}
