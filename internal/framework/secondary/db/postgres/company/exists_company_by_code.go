package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *companyRepo) ExistsCompanyByCode(ctx context.Context, tenantID string, code string, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:ExistsCompanyByCode")
	defer span.End()

	query := `
		SELECT COUNT(*) FROM org_units
		WHERE deleted_at IS NULL AND tenant_id = ? AND category = 'company' AND code = ?
	`
	args := []any{tenantID, code}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID, "code": code}).Msg("Failed to check company code existence")
		return false, err
	}

	return count > 0, nil
}
