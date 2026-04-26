package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalInvoiceCore) ExecuteInvoiceAction(
	ctx context.Context,
	input coreentity.InternalInvoiceActionInput,
) (*coreentity.InternalPaymentAttemptResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:execute_invoice_action:ExecuteInvoiceAction")
	defer span.End()

	switch strings.TrimSpace(input.Action) {
	case coreentity.ActionCreateDOKUPaymentAttempt:
		return c.createDOKUAttempt(ctx, input)
	case coreentity.ActionCreateCustomPaymentAttempt:
		return c.createManualAttempt(ctx, input)
	default:
		return nil, errmsg.NewCustomErrors(400).SetMessage("Unsupported invoice action")
	}
}
