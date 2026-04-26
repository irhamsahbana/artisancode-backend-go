package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *rbacRepo) RemoveRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:RemoveRole")
	defer span.End()

	query := `
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.UserID, data.RoleID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to remove user role")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to get remove user role rows affected")
		return err
	}
	if rowsAffected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("User role not found when removing")
		return errmsg.NewCustomErrors(404).SetMessage("User role not found")
	}
	return nil
}
