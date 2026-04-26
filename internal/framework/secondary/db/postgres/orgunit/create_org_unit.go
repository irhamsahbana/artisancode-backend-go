package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *orgUnitRepo) CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:CreateOrgUnit")
	defer span.End()

	query := `
		INSERT INTO org_units (tenant_id, code, name, parent_id, category)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`

	var id string
	err := r.db.GetContext(
		ctx,
		&id,
		r.db.Rebind(query),
		data.TenantID,
		data.Code,
		data.Name,
		data.ParentID,
		data.Category,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create org unit")
		return nil, err
	}

	data.ID = id
	return &data, nil
}
