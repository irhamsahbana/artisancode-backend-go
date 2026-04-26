package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *orgUnitRepo) UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:UpdateOrgUnit")
	defer span.End()

	query := `
		UPDATE org_units
		SET code = ?, name = ?, parent_id = ?, category = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(query),
		data.Code,
		data.Name,
		data.ParentID,
		data.Category,
		data.ID,
		data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update org unit")
		return err
	}
	return nil
}
