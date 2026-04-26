package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *workLocationRepo) ExistsWorkLocationByName(
	ctx context.Context,
	tenantID string,
	name string,
	excludeID *string,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:worklocation:repo:ExistsWorkLocationByName",
	)
	defer span.End()

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM work_locations
			WHERE tenant_id = ? AND name = ? AND deleted_at IS NULL
		`
	args := []any{tenantID, name}

	if excludeID != nil {
		query += ` AND id != ?`
		args = append(args, *excludeID)
	}

	query += `)`

	err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("name", name).Msg("Failed to check work location name existence")
		return false, err
	}

	return exists, nil
}
