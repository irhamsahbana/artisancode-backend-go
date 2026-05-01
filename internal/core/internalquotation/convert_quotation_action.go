package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (c *internalQuotationCore) executeConvertQuotationAction(
	ctx context.Context,
	input coreentity.InternalQuotationActionInput,
) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalquotation:convert_quotation_action:executeConvertQuotationAction",
	)
	defer span.End()

	var bundle coreentity.InternalCommerceBundle
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		quote, err := c.quotationRepo.GetQuotation(txCtx, input.UserCtx.TenantID, input.ID)
		if err != nil {
			return err
		}
		if quote.Status != coreentity.QuotationStatusApproved {
			log.Ctx(txCtx).
				Warn().
				Str("quotation_id", input.ID).
				Str("current_status", string(quote.Status)).
				Msg("Invalid quotation status for conversion")
			return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageOnlyApprovedQuotationCanBeConverted)
		}
		order, err := c.orderRepo.CreateOrder(txCtx, coreentity.InternalOrder{
			TenantID:                 quote.TenantID,
			CompanyID:                quote.CompanyID,
			Status:                   coreentity.OrderStatusPendingPayment,
			CurrencyCode:             quote.CurrencyCode,
			SubtotalAmount:           quote.SubtotalAmount,
			DiscountAmount:           quote.DiscountAmount,
			TaxAmount:                quote.TaxAmount,
			TotalAmount:              quote.TotalAmount,
			InternalProductID:        quote.InternalProductID,
			InternalProductPricingID: quote.InternalProductPricingID,
			PricingSnapshot:          quote.PricingSnapshot,
			SourceType:               coreentity.OrderSourceTypeQuotation,
			SourceReferenceID:        &quote.ID,
			Metadata:                 quote.Metadata,
		})
		if err != nil {
			return err
		}
		invoice, err := c.invoiceRepo.CreateInvoice(txCtx, coreentity.InternalInvoice{
			InternalOrderID: order.ID, Status: coreentity.InvoiceStatusOpen, CurrencyCode: order.CurrencyCode,
			Amount: order.TotalAmount, AmountPaid: decimal.Zero, AmountOutstanding: order.TotalAmount, DueAt: input.InvoiceDueAt,
			Metadata: map[string]any{},
		})
		if err != nil {
			return err
		}
		converted, err := c.quotationRepo.MarkQuotationConverted(txCtx, quote.TenantID, quote.ID, order.ID)
		if err != nil {
			return err
		}
		bundle.Quotation = converted
		bundle.Order = order
		bundle.Invoice = invoice
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}
