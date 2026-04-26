package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) ExistsInternalUserByEmail(ctx context.Context, email, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaluser:repo:ExistsInternalUserByEmail",
	)
	defer span.End()

	query := `SELECT COUNT(*) FROM internal_users WHERE LOWER(email) = LOWER(?) AND deleted_at IS NULL`
	args := []any{email}
	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"email": email, "exclude_id": excludeID}).
			Msg("Failed to check internal user email")
		return false, err
	}
	return count > 0, nil
}
