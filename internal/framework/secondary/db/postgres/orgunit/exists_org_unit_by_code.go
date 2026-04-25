package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *orgUnitRepo) ExistsOrgUnitByCode(ctx context.Context, tenantID string, code string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:ExistsOrgUnitByCode")
	defer span.End()

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM org_units
			WHERE tenant_id = ? AND code = ? AND deleted_at IS NULL
		)
	`

	err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), tenantID, code)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID, "code": code}).
			Msg("Failed to check org unit code existence")
		return false, err
	}

	return exists, nil
}
