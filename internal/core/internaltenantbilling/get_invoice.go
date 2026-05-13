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

func (c *internalTenantBillingCore) GetInvoice(
	ctx context.Context,
	id string,
) (*coreentity.TenantBillingInvoiceDetail, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_invoice:GetInvoice")
	defer span.End()

	userCtx := common.GetUserContext(ctx)

	detail, err := c.billingRepo.GetInvoice(ctx, userCtx.TenantID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageInvoiceNotFound))
		}
		return nil, err
	}

	return detail, nil
}
