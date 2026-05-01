package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
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
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, input).Msg("Tenant is required to create quotation")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageTenantIsRequired)
	}
	input.CurrencyCode = strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if input.TotalAmount.LessThan(decimal.Zero) || input.SubtotalAmount.LessThan(decimal.Zero) ||
		input.DiscountAmount.LessThan(decimal.Zero) ||
		input.TaxAmount.LessThan(decimal.Zero) {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, input).Msg("Quotation amounts must not be negative")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageAmountsMustNotBeNegative)
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
