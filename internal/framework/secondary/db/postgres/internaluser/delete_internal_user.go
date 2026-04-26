package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) DeleteInternalUser(ctx context.Context, filter coreentity.InternalUserDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:DeleteInternalUser")
	defer span.End()

	result, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(`UPDATE internal_users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`),
		filter.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal user")
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal user not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage("Internal user not found")
	}
	return nil
}
