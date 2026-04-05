package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) ExistsActiveUserByEmailExcludeUser(ctx context.Context, email, excludeUserID string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.ExistsActiveUserByEmailExcludeUser")
	defer span.End()

	var existing string
	query := `
		SELECT id
		FROM users
		WHERE email = ? AND id <> ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &existing, r.db.Rebind(query), email, excludeUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"email":           email,
			"exclude_user_id": excludeUserID,
		}).Msg("Failed to check user existence by email")
		return false, err
	}

	return existing != "", nil
}
