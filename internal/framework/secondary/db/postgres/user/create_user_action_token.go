package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) CreateUserActionToken(ctx context.Context, token coreentity.UserActionToken) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:create_user_action_token:CreateUserActionToken",
	)
	defer span.End()

	query := `
		INSERT INTO user_action_tokens (
			user_id, purpose, token_hash, expires_at
		) VALUES (?, ?, ?, ?)
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(ctx, exec.Rebind(query), token.UserID, token.Purpose, token.TokenHash, token.ExpiresAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, token).Msg("Failed to create user action token")
		return err
	}

	return nil
}
