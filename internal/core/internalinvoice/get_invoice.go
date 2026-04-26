package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalInvoiceCore) GetInvoice(ctx context.Context, id string) (*coreentity.InternalInvoiceDetail, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:get_invoice:GetInvoice")
	defer span.End()

	tenantID := common.GetUserContext(ctx).TenantID
	invoice, err := c.invoiceRepo.GetInvoice(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	order, err := c.orderRepo.GetOrderByInvoiceID(ctx, tenantID, invoice.ID)
	if err != nil {
		return nil, err
	}
	var quote *coreentity.InternalQuotation
	if order.SourceType == coreentity.OrderSourceTypeQuotation && order.SourceReferenceID != nil {
		quote, _ = c.quotationRepo.GetQuotation(ctx, tenantID, *order.SourceReferenceID)
	}
	latest, _ := c.invoiceRepo.GetLatestPaymentAttempt(ctx, tenantID, invoice.ID)
	return &coreentity.InternalInvoiceDetail{
		Invoice:              *invoice,
		Order:                *order,
		Quotation:            quote,
		LatestPaymentAttempt: latest,
	}, nil
}
