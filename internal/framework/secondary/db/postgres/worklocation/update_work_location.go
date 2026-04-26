package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *workLocationRepo) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:worklocation:repo:UpdateWorkLocation")
	defer span.End()

	query := `
		UPDATE work_locations
		SET org_unit_id = ?, name = ?, address = ?, latitude = ?, longitude = ?, radius_meters = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(query),
		data.OrgUnitID,
		data.Name,
		data.Address,
		data.Latitude,
		data.Longitude,
		data.RadiusMeters,
		data.ID,
		data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update work location")
		return err
	}
	return nil
}
