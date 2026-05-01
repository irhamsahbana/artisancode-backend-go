package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) UpdateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProductPricing",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_product_pricings
		SET
			internal_product_id = ?,
			code = ?,
			name = ?,
			description = ?,
			status = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.InternalProductID,
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product pricing")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product pricing not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
	}

	return nil
}
