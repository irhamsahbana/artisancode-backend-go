package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

func (c *internalQuotationCore) CreateQuotation(
	ctx context.Context,
	input coreentity.CreateInternalQuotationInput,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalquotation:create_quotation:CreateQuotation")
	defer span.End()

	tenantID := strings.TrimSpace(input.UserCtx.TenantID)
	if tenantID == "" {
		tenantID = strings.TrimSpace(input.TenantID)
	}
	if tenantID == "" {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Tenant is required")
	}
	input.CurrencyCode = strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if input.TotalAmount.LessThan(decimal.Zero) || input.SubtotalAmount.LessThan(decimal.Zero) ||
		input.DiscountAmount.LessThan(decimal.Zero) ||
		input.TaxAmount.LessThan(decimal.Zero) {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Amounts must not be negative")
	}
	pricingSnapshot, _, err := c.orderRepo.GetPricingSnapshot(
		ctx,
		input.InternalProductID,
		input.InternalProductPricingID,
		input.CurrencyCode,
	)
	if err != nil {
		return nil, err
	}
	return c.quotationRepo.CreateQuotation(ctx, coreentity.InternalQuotation{
		TenantID:                 tenantID,
		CompanyID:                input.UserCtx.CompanyID,
		Status:                   coreentity.QuotationStatusDraft,
		CurrencyCode:             input.CurrencyCode,
		SubtotalAmount:           input.SubtotalAmount,
		DiscountAmount:           input.DiscountAmount,
		TaxAmount:                input.TaxAmount,
		TotalAmount:              input.TotalAmount,
		InternalProductID:        input.InternalProductID,
		InternalProductPricingID: input.InternalProductPricingID,
		PricingSnapshot:          pricingSnapshot,
		QuoteSnapshot:            input.QuoteSnapshot,
		ExpiresAt:                input.ExpiresAt,
		Metadata:                 input.Metadata,
	})
}
