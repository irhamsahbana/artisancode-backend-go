package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalTenantBillingCore) CreatePaymentReceipt(
	ctx context.Context,
	data coreentity.InternalPaymentReceipt,
) (*coreentity.InternalBillingPaymentReceiptActionResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:create_payment_receipt:CreatePaymentReceipt")
	defer span.End()

	data.Status = coreentity.PaymentReceiptStatusPendingVerification

	created, err := c.billingRepo.CreatePaymentReceipt(ctx, data)
	if err != nil {
		return nil, err
	}

	return &coreentity.InternalBillingPaymentReceiptActionResult{
		Receipt: *created,
		Status:  created.Status,
	}, nil
}
