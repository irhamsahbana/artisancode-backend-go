package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) MarkUserActionTokenUsed(ctx context.Context, tokenID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:mark_user_action_token_used:MarkUserActionTokenUsed")
	defer span.End()

	query := `
		UPDATE user_action_tokens
		SET used_at = NOW(), updated_at = NOW()
		WHERE id = ? AND used_at IS NULL AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(ctx, exec.Rebind(query), tokenID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"token_id": tokenID,
		}).Msg("Failed to mark user action token as used")
		return err
	}

	return nil
}
