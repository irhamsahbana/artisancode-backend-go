package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProduct",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_products
		SET
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
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
	}

	return nil
}
