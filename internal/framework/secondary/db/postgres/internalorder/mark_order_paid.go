package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalOrderRepo) MarkOrderPaid(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalorder:mark_order_paid:MarkOrderPaid",
	)
	defer span.End()

	query := `
		UPDATE internal_orders
		SET
			status = 'paid',
			updated_at = NOW()
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), id, tenantID); err != nil {
		return nil, err
	}
	return r.GetOrder(ctx, tenantID, id)
}
