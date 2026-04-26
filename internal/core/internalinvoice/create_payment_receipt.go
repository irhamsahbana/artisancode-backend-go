package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalInvoiceCore) CreatePaymentReceipt(
	ctx context.Context,
	receipt coreentity.InternalPaymentReceipt,
) (*coreentity.InternalPaymentReceiptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:create_payment_receipt:CreatePaymentReceipt")
	defer span.End()

	tenantID := receipt.UserCtx.TenantID
	if receipt.Status == "" {
		receipt.Status = coreentity.PaymentReceiptStatusPendingVerification
	}
	receipt.SourceType = coreentity.PaymentProviderManual
	if receipt.Status == coreentity.PaymentReceiptStatusAccepted {
		now := time.Now().UTC().Format(time.RFC3339)
		receipt.VerifiedAt = &now
		receipt.VerifiedByUserID = &receipt.UserCtx.UserID
	}
	var result coreentity.InternalPaymentReceiptResult
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		invoice, err := c.invoiceRepo.GetInvoice(txCtx, tenantID, receipt.InternalInvoiceID)
		if err != nil {
			return err
		}
		if receipt.AmountReceived.LessThan(invoice.AmountOutstanding) &&
			receipt.Status == coreentity.PaymentReceiptStatusAccepted {
			return errmsg.NewCustomErrors(400).SetMessage("Payment amount does not match invoice outstanding amount")
		}
		created, err := c.invoiceRepo.CreatePaymentReceipt(txCtx, receipt)
		if err != nil {
			return err
		}
		order, err := c.orderRepo.GetOrderByInvoiceID(txCtx, tenantID, invoice.ID)
		if err != nil {
			return err
		}
		if receipt.Status == coreentity.PaymentReceiptStatusAccepted {
			invoice, err = c.invoiceRepo.MarkInvoicePaid(txCtx, tenantID, invoice.ID, receipt.AmountReceived.String())
			if err != nil {
				return err
			}
			order, err = c.orderRepo.MarkOrderPaid(txCtx, tenantID, order.ID)
			if err != nil {
				return err
			}
		}
		result = coreentity.InternalPaymentReceiptResult{Receipt: *created, Invoice: *invoice, Order: *order}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
