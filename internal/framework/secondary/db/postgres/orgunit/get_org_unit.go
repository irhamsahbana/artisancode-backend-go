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

func (r *orgUnitRepo) GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:orgunit:repo:GetOrgUnit")
	defer span.End()

	var data struct {
		ID       string  `db:"id"`
		TenantID string  `db:"tenant_id"`
		Code     string  `db:"code"`
		Name     string  `db:"name"`
		ParentID *string `db:"parent_id"`
		Category string  `db:"category"`
	}

	query := `
		SELECT id, tenant_id, code, name, parent_id, category
		FROM org_units
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageOrgUnitNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageOrgUnitNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get org unit")
		return nil, err
	}

	result := &coreentity.OrgUnit{
		ID:       data.ID,
		TenantID: data.TenantID,
		Code:     data.Code,
		Name:     data.Name,
		ParentID: data.ParentID,
		Category: data.Category,
	}
	return result, nil
}
