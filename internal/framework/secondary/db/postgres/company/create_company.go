package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *companyRepo) CreateCompany(ctx context.Context, data coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:CreateCompany")
	defer span.End()

	configJSON, err := json.Marshal(data.Config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to marshal config")
		return nil, err
	}

	query := `
		INSERT INTO org_units (tenant_id, category, name, code, config)
		VALUES (?, 'company', ?, ?, ?)
		RETURNING id
	`

	var id string
	err = r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.Name, data.Code, configJSON,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create company")
		return nil, err
	}

	data.ID = id
	return &data, nil
}
