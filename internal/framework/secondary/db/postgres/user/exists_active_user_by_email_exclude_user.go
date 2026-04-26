package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) ExistsActiveUserByEmailAndTenantExcludeUser(
	ctx context.Context,
	email, tenantID, excludeUserID string,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:exists_active_user_by_email_exclude_user:ExistsActiveUserByEmailAndTenantExcludeUser",
	)
	defer span.End()

	var existing string
	query := `
		SELECT id
		FROM users
		WHERE email = ? AND tenant_id = ? AND id <> ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &existing, exec.Rebind(query), email, tenantID, excludeUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"email":           email,
			"tenant_id":       tenantID,
			"exclude_user_id": excludeUserID,
		}).Msg("Failed to check user existence by email")
		return false, err
	}

	return existing != "", nil
}
