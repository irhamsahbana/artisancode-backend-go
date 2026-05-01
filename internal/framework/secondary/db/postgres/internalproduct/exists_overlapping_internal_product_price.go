package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) ExistsOverlappingInternalProductPrice(
	ctx context.Context,
	filter coreentity.InternalProductPriceOverlapFilter,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsOverlappingInternalProductPrice",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_product_prices
		WHERE internal_product_pricing_id = ?
		  AND currency_code = ?
		  AND deleted_at IS NULL
		  AND started_at < COALESCE(?, 'infinity'::timestamptz)
		  AND COALESCE(ended_at, 'infinity'::timestamptz) > ?
	`
	args := []any{filter.InternalProductPricingID, filter.CurrencyCode, filter.EndedAt, filter.StartedAt}
	if filter.ExcludeID != "" {
		query += ` AND id != ?`
		args = append(args, filter.ExcludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, filter).
			Msg("Failed to check overlapping internal product price")
		return false, err
	}

	return count > 0, nil
}
