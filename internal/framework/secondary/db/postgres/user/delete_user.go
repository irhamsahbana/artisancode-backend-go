package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) DeleteUser(ctx context.Context, filter coreentity.UserDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:delete_user:DeleteUser")
	defer span.End()

	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	result, err := exec.ExecContext(ctx, exec.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get delete user rows affected")
		return err
	}
	if rowsAffected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("User not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageUserNotFound)
	}

	return nil
}
