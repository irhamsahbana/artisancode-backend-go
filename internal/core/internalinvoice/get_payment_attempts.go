package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalInvoiceCore) GetPaymentAttempts(
	ctx context.Context,
	invoiceID string,
) ([]coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalinvoice:get_payment_attempts:GetPaymentAttempts")
	defer span.End()

	return c.invoiceRepo.GetPaymentAttempts(ctx, common.GetUserContext(ctx).TenantID, invoiceID)
}
