package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalInvoiceCore) ExecutePaymentAttemptAction(
	ctx context.Context,
	input coreentity.InternalPaymentAttemptActionInput,
) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalinvoice:execute_payment_attempt_action:ExecutePaymentAttemptAction",
	)
	defer span.End()

	if input.Action != coreentity.ActionRetryPayment {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageUnsupportedPaymentAttemptAction)
	}
	attempt, err := c.invoiceRepo.GetPaymentAttempt(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != coreentity.PaymentAttemptStatusFailed &&
		attempt.Status != coreentity.PaymentAttemptStatusExpired &&
		attempt.Status != coreentity.PaymentAttemptStatusCancelled {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePaymentAttemptCannotBeRetried)
	}
	return c.ExecuteInvoiceAction(ctx, coreentity.InternalInvoiceActionInput{
		UserCtx: input.UserCtx, ID: attempt.InternalInvoiceID, Action: actionForProvider(input.Provider, attempt.Provider),
		Provider: input.Provider, PaymentMethodType: input.PaymentMethodType, PaymentChannelCode: input.PaymentChannelCode,
		CallbackURL: input.CallbackURL, CustomerName: input.CustomerName, CustomerEmail: input.CustomerEmail, CustomerPhone: input.CustomerPhone,
		RequestID: input.RequestID,
	})
}
