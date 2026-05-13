package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) GetPaymentAttempts(
	ctx context.Context,
	invoiceID string,
) ([]coreentity.TenantBillingPaymentAttemptView, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_payment_attempts:GetPaymentAttempts")
	defer span.End()

	userCtx := common.GetUserContext(ctx)

	_, err := c.billingRepo.GetInvoiceRaw(ctx, userCtx.TenantID, invoiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageInvoiceNotFound))
		}
		return nil, err
	}

	return c.billingRepo.GetPaymentAttemptsByInvoice(ctx, userCtx.TenantID, invoiceID)
}
