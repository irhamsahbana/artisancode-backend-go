package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/shopspring/decimal"
)

func (c *internalTenantBillingCore) ExecutePaymentAttemptAction(
	ctx context.Context,
	input coreentity.TenantBillingPaymentAttemptActionInput,
) (*coreentity.TenantBillingPaymentAttemptActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:execute_payment_attempt_action:ExecutePaymentAttemptAction")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	input.UserCtx = userCtx

	switch input.Action {
	case coreentity.TenantBillingPaymentAttemptActionRetry:
		return c.retryPaymentAttempt(ctx, input)
	default:
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedPaymentAttemptAction))
	}
}

func (c *internalTenantBillingCore) retryPaymentAttempt(
	ctx context.Context,
	input coreentity.TenantBillingPaymentAttemptActionInput,
) (*coreentity.TenantBillingPaymentAttemptActionResult, error) {
	attempt, err := c.billingRepo.GetPaymentAttemptByID(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessagePaymentAttemptNotFound))
		}
		return nil, err
	}

	if attempt.Status != coreentity.InternalTenantPaymentAttemptStatusFailed &&
		attempt.Status != coreentity.InternalTenantPaymentAttemptStatusExpired {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessagePaymentAttemptCannotBeRetried))
	}

	invoice, err := c.billingRepo.GetInvoiceRaw(ctx, input.UserCtx.TenantID, attempt.InternalTenantInvoiceID)
	if err != nil {
		return nil, err
	}

	amount := decimal.RequireFromString(invoice.Amount)
	multiplier := decimal.NewFromInt(1)
	if invoice.CurrencyCode != "" {
		multiplier = decimal.NewFromInt(1)
	}
	dokuAmount := amount.Mul(multiplier).Round(0).IntPart()
	if dokuAmount <= 0 {
		dokuAmount = amount.Mul(decimal.NewFromInt(100)).Round(0).IntPart()
	}

	dokuResp, err := c.doku.CreatePayment(ctx, restentity.DokuCreatePaymentRequest{
		InvoiceNumber:   invoice.InvoiceNumber,
		Amount:          dokuAmount,
		Currency:        invoice.CurrencyCode,
		CustomerName:    checkoutCustomerName(input.UserCtx),
		CustomerEmail:   checkoutCustomerEmail(input.UserCtx),
		LineItems:       []restentity.DokuLineItem{{Name: invoice.InvoiceNumber, Price: dokuAmount, Quantity: 1}},
	})
	if err != nil {
		return nil, err
	}

	attempt.ProviderRequestID = dokuResp.RequestID
	attempt.ProviderPaymentURL = dokuResp.PaymentURL
	attempt.ProviderPayloadSnapshot = map[string]any{"doku_response": dokuResp}
	attempt.Status = coreentity.InternalTenantPaymentAttemptStatusPending
	attempt, err = c.billingRepo.UpdatePaymentAttemptGateway(ctx, *attempt)
	if err != nil {
		return nil, err
	}

	return &coreentity.TenantBillingPaymentAttemptActionResult{
		PaymentAttemptID: attempt.ID,
		Status:           attempt.Status,
		PaymentURL:       attempt.ProviderPaymentURL,
	}, nil
}
