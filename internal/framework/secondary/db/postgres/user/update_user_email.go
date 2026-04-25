package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) UpdateUserEmail(ctx context.Context, userID, tenantID, email string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:update_user_email:UpdateUserEmail")
	defer span.End()

	query := `
		UPDATE users
		SET email = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	_, err := exec.ExecContext(ctx, exec.Rebind(query), email, userID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"user_id":   userID,
			"tenant_id": tenantID,
			"email":     email,
		}).Msg("Failed to update user email")
		return err
	}

	return nil
}
