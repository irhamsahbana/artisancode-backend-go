package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) CreateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) (*coreentity.InternalProductPrice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:CreateInternalProductPrice",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO internal_product_prices (
			internal_product_pricing_id,
			currency_code,
			amount,
			started_at,
			ended_at,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id
	`
	if err := r.db.GetContext(ctx, &data.ID, r.db.Rebind(query),
		data.InternalProductPricingID,
		data.CurrencyCode,
		data.Amount,
		data.StartedAt,
		data.EndedAt,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal product price")
		return nil, err
	}

	return &data, nil
}
