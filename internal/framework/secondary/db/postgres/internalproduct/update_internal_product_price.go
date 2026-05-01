package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) UpdateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProductPrice",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_product_prices
		SET
			internal_product_pricing_id = ?,
			currency_code = ?,
			amount = ?,
			started_at = ?,
			ended_at = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.InternalProductPricingID,
		data.CurrencyCode,
		data.Amount,
		data.StartedAt,
		data.EndedAt,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product price")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product price not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPriceNotFound)
	}

	return nil
}
