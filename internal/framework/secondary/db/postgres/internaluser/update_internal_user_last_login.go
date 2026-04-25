package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) UpdateInternalUserLastLogin(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:UpdateInternalUserLastLogin")
	defer span.End()

	_, err := r.db.ExecContext(ctx, r.db.Rebind(`UPDATE internal_users SET last_login_at = NOW(), updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`), id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, id).Msg("Failed to update internal user last login")
		return err
	}
	return nil
}
