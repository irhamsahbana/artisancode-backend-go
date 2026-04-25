package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) DeleteUserActionTokensByPurpose(ctx context.Context, userID, purpose string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:delete_user_action_tokens_by_purpose:DeleteUserActionTokensByPurpose")
	defer span.End()

	query := `
		UPDATE user_action_tokens
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE user_id = ? AND purpose = ? AND used_at IS NULL AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(ctx, exec.Rebind(query), userID, purpose)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id": userID,
			"purpose": purpose,
		}).Msg("Failed to invalidate user action tokens")
		return err
	}

	return nil
}
