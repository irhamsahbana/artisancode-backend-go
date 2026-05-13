package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetSubscriptionsInGracePastDue(
	ctx context.Context,
	limit int,
) ([]coreentity.InternalTenantSubscription, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_subscriptions_in_grace_past_due:GetSubscriptionsInGracePastDue",
	)
	defer span.End()

	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT
			id,
			internal_billing_account_id,
			tenant_id,
			product_family,
			status,
			internal_product_id,
			internal_product_pricing_id,
			plan_snapshot,
			add_on_snapshots,
			current_period_started_at,
			current_period_ended_at,
			grace_ended_at,
			metadata
		FROM internal_tenant_subscriptions
		WHERE status = 'grace_period'
			AND grace_ended_at IS NOT NULL
			AND grace_ended_at <= NOW()
			AND deleted_at IS NULL
		ORDER BY grace_ended_at ASC
		LIMIT ?
	`

	items := make([]coreentity.InternalTenantSubscription, 0)
	if err := r.exec(ctx).SelectContext(ctx, &items, r.exec(ctx).Rebind(query), limit); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get subscriptions in grace past due")
		return nil, err
	}

	return items, nil
}
