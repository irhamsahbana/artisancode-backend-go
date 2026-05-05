package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

func (c *internalInvoiceCore) createDOKUAttempt(
	ctx context.Context,
	input coreentity.InternalInvoiceActionInput,
) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:payment_attempt_helpers:createDOKUAttempt")
	defer span.End()

	if c.doku == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageDokuClientIsNotConfigured)
	}
	invoice, err := c.invoiceRepo.GetInvoice(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if invoice.Status != coreentity.InvoiceStatusOpen && invoice.Status != coreentity.InvoiceStatusPartiallyPaid {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvoiceIsNotPayable)
	}

	currency, err := c.currencyRepo.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{
		Code: invoice.CurrencyCode,
	})
	if err != nil {
		return nil, err
	}

	amount := invoice.AmountOutstanding
	multiplier := decimal.NewFromInt(1)
	if currency.DecimalPlaces > 0 {
		multiplier = decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(currency.DecimalPlaces)))
	}
	dokuAmount := amount.Mul(multiplier).Round(0).IntPart()

	attempt, err := c.invoiceRepo.CreatePaymentAttempt(ctx, coreentity.InternalPaymentAttempt{
		InternalInvoiceID: invoice.ID, Provider: coreentity.PaymentProviderDOKU, PaymentMethodType: "gateway",
		ProviderReference: invoice.InvoiceNumber, Status: coreentity.PaymentAttemptStatusInitiated,
		RequestedAmount: invoice.AmountOutstanding, PaidAmount: decimal.Zero, Metadata: map[string]any{},
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.doku.CreatePayment(ctx, restentity.DokuCreatePaymentRequest{
		InvoiceNumber: invoice.InvoiceNumber, Amount: dokuAmount,
		CustomerName: input.CustomerName, CustomerEmail: input.CustomerEmail, CustomerPhone: input.CustomerPhone,
		CallbackURL: input.CallbackURL,
	})
	if err != nil {
		return nil, err
	}
	attempt.ProviderRequestID = resp.RequestID
	attempt.ProviderPaymentURL = resp.PaymentURL
	attempt.ProviderPayloadSnapshot = map[string]any{"doku_response": resp}
	attempt.Status = coreentity.PaymentAttemptStatusPending
	attempt, err = c.invoiceRepo.UpdatePaymentAttemptGateway(ctx, *attempt)
	if err != nil {
		return nil, err
	}
	return &coreentity.InternalPaymentAttemptResult{PaymentAttempt: *attempt}, nil
}

func (c *internalInvoiceCore) createManualAttempt(
	ctx context.Context,
	input coreentity.InternalInvoiceActionInput,
) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:payment_attempt_helpers:createManualAttempt")
	defer span.End()

	invoice, err := c.invoiceRepo.GetInvoice(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if invoice.Status != coreentity.InvoiceStatusOpen && invoice.Status != coreentity.InvoiceStatusPartiallyPaid {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInvoiceIsNotPayable)
	}
	provider := input.Provider
	if provider == "" {
		provider = coreentity.PaymentProviderManual
	}
	method := input.PaymentMethodType
	if method == "" {
		method = "bank_transfer"
	}
	channel := input.PaymentChannelCode
	if channel == "" {
		channel = "manual_bank_transfer"
	}
	attempt, err := c.invoiceRepo.CreatePaymentAttempt(ctx, coreentity.InternalPaymentAttempt{
		InternalInvoiceID: invoice.ID, Provider: provider, PaymentMethodType: method, PaymentChannelCode: channel,
		ProviderReference: invoice.InvoiceNumber, Status: coreentity.PaymentAttemptStatusPending,
		RequestedAmount: invoice.AmountOutstanding, PaidAmount: decimal.Zero, Metadata: map[string]any{},
	})
	if err != nil {
		return nil, err
	}
	return &coreentity.InternalPaymentAttemptResult{
		PaymentAttempt: *attempt,
		Instruction: &coreentity.InternalManualPaymentInstruction{
			Title: "Manual transfer", ReferenceNumber: invoice.InvoiceNumber, BankName: "BCA",
			AccountNumber: "1234567890", AccountName: "PT Artisan Code", Notes: "Use invoice number as transfer note",
		},
	}, nil
}

func actionForProvider(requestProvider, previousProvider string) string {
	provider := requestProvider
	if provider == "" {
		provider = previousProvider
	}
	if provider == coreentity.PaymentProviderDOKU {
		return coreentity.ActionCreateDOKUPaymentAttempt
	}
	return coreentity.ActionCreateCustomPaymentAttempt
}
