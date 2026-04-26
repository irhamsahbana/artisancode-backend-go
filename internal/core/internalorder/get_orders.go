package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalOrderCore) GetOrders(
	ctx context.Context,
	filter coreentity.InternalCommerceListFilter,
) ([]coreentity.InternalOrder, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalorder:get_orders:GetOrders")
	defer span.End()

	if filter.TenantID == "" {
		filter.TenantID = common.GetUserContext(ctx).TenantID
	}
	return c.orderRepo.GetOrders(ctx, filter)
}
