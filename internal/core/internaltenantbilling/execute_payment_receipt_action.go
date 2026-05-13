package core

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) ExecutePaymentReceiptAction(
	ctx context.Context,
	input coreentity.InternalBillingPaymentReceiptActionInput,
) (*coreentity.InternalBillingPaymentReceiptActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:execute_payment_receipt_action:ExecutePaymentReceiptAction")
	defer span.End()

	receipt, err := c.billingRepo.GetPaymentReceipt(ctx, input.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("payment receipt tidak ditemukan"))
		}
		return nil, err
	}

	if receipt.Status != coreentity.PaymentReceiptStatusPendingVerification {
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("payment receipt sudah diproses sebelumnya"))
	}

	now := time.Now()
	userCtx := common.GetUserContext(ctx)

	switch input.Action {
	case coreentity.PaymentReceiptActionAccept:
		receipt.Status = coreentity.PaymentReceiptStatusAccepted
		receipt.VerifiedAt = ptr(now.Format(time.RFC3339))
		receipt.VerifiedByUserID = &userCtx.UserID

	case coreentity.PaymentReceiptActionReject:
		receipt.Status = coreentity.PaymentReceiptStatusRejected
		receipt.VerifiedAt = ptr(now.Format(time.RFC3339))
		receipt.VerifiedByUserID = &userCtx.UserID

	default:
		return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("unsupported action: "+input.Action))
	}

	updated, err := c.billingRepo.UpdatePaymentReceipt(ctx, *receipt)
	if err != nil {
		return nil, err
	}

	return &coreentity.InternalBillingPaymentReceiptActionResult{
		Receipt: *updated,
		Status:  updated.Status,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}
