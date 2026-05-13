package core

import (
	"context"
	"database/sql"
	"errors"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) GetInvoiceByID(
	ctx context.Context,
	id string,
) (*coreentity.InternalTenantInvoice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:get_invoice_by_id:GetInvoiceByID")
	defer span.End()

	item, err := c.billingRepo.GetInvoiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage(errmsg.MessageInvoiceNotFound))
		}
		return nil, err
	}
	return item, nil
}
