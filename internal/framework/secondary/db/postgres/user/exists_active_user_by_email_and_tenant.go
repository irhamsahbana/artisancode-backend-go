package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) ExistsActiveUserByEmailAndTenant(ctx context.Context, email, tenantID string) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:register_owner:ExistsActiveUserByEmailAndTenant",
	)
	defer span.End()

	var existing string
	query := `SELECT id FROM users WHERE email = ? AND tenant_id = ? AND deleted_at IS NULL`
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &existing, exec.Rebind(query), email, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email":     email,
				"tenant_id": tenantID,
			}).Msg("User not found")
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"email":     email,
			"tenant_id": tenantID,
		}).Msg("Failed to check user existence")
		return false, err
	}
	return existing != "", nil
}
