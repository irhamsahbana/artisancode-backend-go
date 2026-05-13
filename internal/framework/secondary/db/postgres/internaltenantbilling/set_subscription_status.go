package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) SetSubscriptionStatus(
	ctx context.Context,
	tenantID string,
	subscriptionID string,
	status string,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:set_subscription_status:SetSubscriptionStatus",
	)
	defer span.End()

	query := `
		UPDATE internal_tenant_subscriptions
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	_, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		status,
		subscriptionID,
		tenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("subscription_id", subscriptionID).Str("status", status).Msg("Failed to set subscription status")
		return err
	}

	return nil
}
