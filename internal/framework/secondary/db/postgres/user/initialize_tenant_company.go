package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) insertDefaultCompany(
	ctx context.Context,
	tx sqlExecutor,
	tenantID string,
	companyName string,
	preferredLanguage string,
) (string, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:insertDefaultCompany",
	)
	defer span.End()

	config := coreentity.DefaultCompanyConfig()
	if preferredLanguage != "" {
		config.PreferredLanguage = preferredLanguage
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "companyName": companyName}).
			Msg("Failed to marshal company config")
		return "", err
	}

	query := `
		INSERT INTO org_units (tenant_id, name, code, category, config)
		VALUES (?, ?, ?, 'company', ?)
		RETURNING id
	`

	var companyID string
	err = tx.GetContext(ctx, &companyID, tx.Rebind(query), tenantID, companyName, companyName, configJSON)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "companyName": companyName}).
			Msg("Failed to insert default company")
		return "", err
	}

	return companyID, nil
}
