package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) ExistsTenantByCode(ctx context.Context, code string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:ExistsTenantByCode")
	defer span.End()

	var existing string
	query := `SELECT id FROM tenants WHERE code = ? AND deleted_at IS NULL`
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &existing, exec.Rebind(query), code)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"code": code}).Msg("Failed to check tenant existence")
		return false, err
	}
	return existing != "", nil
}
