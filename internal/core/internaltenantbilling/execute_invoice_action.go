package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *internalTenantBillingCore) ExecuteInvoiceAction(
	ctx context.Context,
	input coreentity.TenantBillingInvoiceActionInput,
) (*coreentity.TenantBillingInvoiceActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:execute_invoice_action:ExecuteInvoiceAction")
	defer span.End()

	userCtx := common.GetUserContext(ctx)
	input.UserCtx = userCtx

	switch input.Action {
	case coreentity.TenantBillingInvoiceActionCancelCheckout:
		return c.cancelCheckoutOrder(ctx, input)
	default:
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageUnsupportedInvoiceAction))
	}
}

func (c *internalTenantBillingCore) cancelCheckoutOrder(
	ctx context.Context,
	input coreentity.TenantBillingInvoiceActionInput,
) (*coreentity.TenantBillingInvoiceActionResult, error) {
	invoice, err := c.billingRepo.GetInvoiceRaw(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageInvoiceNotFound))
		}
		return nil, err
	}

	if invoice.Status != coreentity.InternalTenantInvoiceStatusPending {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage(errmsg.MessageInvoiceIsNotPayable))
	}

	if err := c.billingRepo.CancelInvoice(ctx, input.UserCtx.TenantID, input.ID); err != nil {
		return nil, err
	}

	if invoice.InvoiceNumber != "" && c.doku != nil {
		if _, err := c.doku.CancelOrder(ctx, invoice.InvoiceNumber); err != nil {
			log.Ctx(ctx).Warn().Err(err).
				Str("invoice_id", input.ID).
				Str("doku_invoice_number", invoice.InvoiceNumber).
				Msg("best-effort DOKU cancel order failed")
		}
	}

	return &coreentity.TenantBillingInvoiceActionResult{
		InvoiceID: input.ID,
		Status:    coreentity.InternalTenantInvoiceStatusCancelled,
	}, nil
}
