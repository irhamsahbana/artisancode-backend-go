package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) UpdateUserPassword(ctx context.Context, userID, tenantID, hashedPassword string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:update_user_password:UpdateUserPassword")
	defer span.End()

	query := `
		UPDATE users
		SET password = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), hashedPassword, userID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id":   userID,
			"tenant_id": tenantID,
		}).Msg("Failed to update user password")
		return err
	}

	return nil
}
