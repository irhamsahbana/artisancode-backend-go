package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalOrderCore) GetOrder(ctx context.Context, id string) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:get_order:GetOrder")
	defer span.End()

	tenantID := common.GetUserContext(ctx).TenantID
	order, err := c.orderRepo.GetOrder(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	invoice, _ := c.invoiceRepo.GetInvoiceByOrderID(ctx, tenantID, id)
	return &coreentity.InternalCommerceBundle{Order: order, Invoice: invoice}, nil
}
