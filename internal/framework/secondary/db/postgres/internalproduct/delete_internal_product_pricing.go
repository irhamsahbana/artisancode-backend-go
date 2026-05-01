package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) DeleteInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:DeleteInternalProductPricing",
	)
	defer span.End()

	query := `
		UPDATE internal_product_pricings
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal product pricing")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product pricing not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
	}

	return nil
}
