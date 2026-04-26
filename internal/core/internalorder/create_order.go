package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

func (c *internalOrderCore) CreateOrder(
	ctx context.Context,
	input coreentity.CreateInternalOrderInput,
) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:create_order:CreateOrder")
	defer span.End()

	uc := input.UserCtx
	tenantID := strings.TrimSpace(uc.TenantID)
	if tenantID == "" {
		tenantID = strings.TrimSpace(input.TenantID)
	}
	if tenantID == "" {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Tenant is required")
	}
	currency := strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	pricingSnapshot, amount, err := c.orderRepo.GetPricingSnapshot(
		ctx,
		input.InternalProductID,
		input.InternalProductPricingID,
		currency,
	)
	if err != nil {
		return nil, err
	}
	total, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}
	var bundle coreentity.InternalCommerceBundle
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		order, err := c.orderRepo.CreateOrder(txCtx, coreentity.InternalOrder{
			TenantID: tenantID, CompanyID: uc.CompanyID, Status: coreentity.OrderStatusPendingPayment, CurrencyCode: currency,
			SubtotalAmount: total, DiscountAmount: decimal.Zero, TaxAmount: decimal.Zero, TotalAmount: total,
			InternalProductID: input.InternalProductID, InternalProductPricingID: input.InternalProductPricingID, PricingSnapshot: pricingSnapshot,
			SourceType: coreentity.OrderSourceTypeStandardPricing, Metadata: input.Metadata,
		})
		if err != nil {
			return err
		}
		invoice, err := c.invoiceRepo.CreateInvoice(txCtx, coreentity.InternalInvoice{
			InternalOrderID: order.ID, Status: coreentity.InvoiceStatusOpen, CurrencyCode: currency,
			Amount: total, AmountPaid: decimal.Zero, AmountOutstanding: total, DueAt: input.InvoiceDueAt, Metadata: map[string]any{},
		})
		if err != nil {
			return err
		}
		bundle.Order = order
		bundle.Invoice = invoice
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}
