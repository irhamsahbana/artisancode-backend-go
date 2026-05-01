package repository

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) UpdateInternalUser(ctx context.Context, data coreentity.InternalUser) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:UpdateInternalUser")
	defer span.End()

	args := []any{data.Email, data.FullName, data.RoleCode, data.Status}
	query := `
		UPDATE internal_users
		SET email = ?, full_name = ?, role_code = ?, status = ?
	`
	if strings.TrimSpace(data.Password) != "" {
		query += `, password_hash = ?`
		args = append(args, data.Password)
	}
	query += `, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	args = append(args, data.ID)

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal user")
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal user not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalUserNotFound)
	}
	return nil
}
