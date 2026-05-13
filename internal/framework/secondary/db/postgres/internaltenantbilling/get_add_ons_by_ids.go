package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetAddOnsByIDs(
	ctx context.Context,
	addOnIDs []string,
) ([]coreentity.TenantBillingAddOn, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_add_ons_by_ids:GetAddOnsByIDs",
	)
	defer span.End()

	if len(addOnIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT
			ip.id,
			ip.code,
			ip.name,
			ip.description,
			ipp.amount::text,
			ipp.currency_code,
			ipp.billing_cycle,
			ipp.id AS "pricing_id",
			ip.metadata AS "raw_metadata"
		FROM internal_products ip
		JOIN internal_product_pricings ipp
			ON ipp.internal_product_id = ip.id
			AND ipp.deleted_at IS NULL
		WHERE ip.id IN (?)
			AND ip.deleted_at IS NULL
			AND ip.status = 'active'
			AND ip.metadata->>'product_type' IS NOT DISTINCT FROM 'add_on'
	`

	reboundQuery, args, err := sqlx.In(query, addOnIDs)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any("add_on_ids", addOnIDs).Msg("Failed to build add-ons query")
		return nil, err
	}
	reboundQuery = r.exec(ctx).Rebind(reboundQuery)

	items := make([]coreentity.TenantBillingAddOn, 0)
	if err := r.exec(ctx).SelectContext(ctx, &items, reboundQuery, args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any("add_on_ids", addOnIDs).Msg("Failed to get add-ons by IDs")
		return nil, err
	}

	return items, nil
}
