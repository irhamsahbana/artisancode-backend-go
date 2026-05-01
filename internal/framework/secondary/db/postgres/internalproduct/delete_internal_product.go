package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) DeleteInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:DeleteInternalProduct",
	)
	defer span.End()

	query := `
		UPDATE internal_products
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal product")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
	}

	return nil
}
