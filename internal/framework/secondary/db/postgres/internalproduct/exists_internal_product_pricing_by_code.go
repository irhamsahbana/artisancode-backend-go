package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) ExistsInternalProductPricingByCode(
	ctx context.Context,
	internalProductID string,
	code string,
	excludeID string,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsInternalProductPricingByCode",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_product_pricings
		WHERE internal_product_id = ? AND code = ? AND deleted_at IS NULL
	`
	args := []any{internalProductID, code}
	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"internal_product_id": internalProductID,
			"code":                code,
			"exclude_id":          excludeID,
		}).Msg("Failed to check internal product pricing code")
		return false, err
	}

	return count > 0, nil
}
