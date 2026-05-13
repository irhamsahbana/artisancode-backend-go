package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetPricingInfo(ctx context.Context, pricingID string) (*coreentity.PricingInfo, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_pricing_info:GetPricingInfo",
	)
	defer span.End()

	query := `
		SELECT
			ipp.amount::text,
			ipp.currency_code,
			ipp.billing_cycle
		FROM internal_product_pricings ipp
		WHERE ipp.id = ?
			AND ipp.deleted_at IS NULL
	`

	var info coreentity.PricingInfo
	if err := r.exec(ctx).GetContext(ctx, &info, r.exec(ctx).Rebind(query), pricingID); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("pricing_id", pricingID).Msg("Failed to get pricing info")
		return nil, err
	}

	return &info, nil
}
