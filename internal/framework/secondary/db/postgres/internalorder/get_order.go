package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalOrderRepo) GetOrder(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalorder:get_order:GetOrder")
	defer span.End()

	return r.getOrderByWhere(ctx, tenantID, "o.id = ?", id)
}
