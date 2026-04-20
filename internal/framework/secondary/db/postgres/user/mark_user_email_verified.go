package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) MarkUserEmailVerified(ctx context.Context, userID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:mark_user_email_verified:MarkUserEmailVerified")
	defer span.End()

	query := `
		UPDATE users
		SET email_verified_at = COALESCE(email_verified_at, NOW()), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), userID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id": userID,
		}).Msg("Failed to mark user email as verified")
		return err
	}

	return nil
}
