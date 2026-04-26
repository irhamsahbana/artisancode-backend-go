package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *attendanceRepo) GetWorkLocation(
	ctx context.Context,
	filter coreentity.WorkLocation,
) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:attendance:get_work_location:GetWorkLocation",
	)
	defer span.End()

	var data struct {
		ID           string   `db:"id"`
		TenantID     string   `db:"tenant_id"`
		RadiusMeters *int     `db:"radius_meters"`
		Latitude     *float64 `db:"latitude"`
		Longitude    *float64 `db:"longitude"`
	}

	query := `
		SELECT id, tenant_id, radius_meters, latitude, longitude
		FROM work_locations
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Work location not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage("Work location not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get work location")
		return nil, err
	}

	return &coreentity.WorkLocation{
		ID:           data.ID,
		TenantID:     data.TenantID,
		RadiusMeters: data.RadiusMeters,
		Latitude:     data.Latitude,
		Longitude:    data.Longitude,
	}, nil
}
