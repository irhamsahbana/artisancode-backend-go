package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) ExistsInternalProductByCode(ctx context.Context, code, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsInternalProductByCode",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_products
		WHERE code = ? AND deleted_at IS NULL
	`
	args := []any{code}
	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"code":       code,
			"exclude_id": excludeID,
		}).Msg("Failed to check internal product code")
		return false, err
	}

	return count > 0, nil
}
