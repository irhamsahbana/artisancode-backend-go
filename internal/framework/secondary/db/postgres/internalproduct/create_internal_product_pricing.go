package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) CreateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:CreateInternalProductPricing",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO internal_product_pricings (
			internal_product_id,
			code,
			name,
			description,
			status,
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
		data.InternalProductID,
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal product pricing")
		return nil, err
	}

	return &data, nil
}
