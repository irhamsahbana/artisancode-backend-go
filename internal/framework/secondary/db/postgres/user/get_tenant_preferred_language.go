package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetTenantPreferredLanguage(ctx context.Context, tenantID string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:get_tenant_preferred_language:GetTenantPreferredLanguage")
	defer span.End()

	query := `
		SELECT config
		FROM org_units
		WHERE tenant_id = ? AND category = 'company' AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`

	var configJSON []byte
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &configJSON, exec.Rebind(query), tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
		}).Msg("Failed to get tenant preferred language config")
		return "", err
	}

	var config coreentity.CompanyConfig
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
		}).Msg("Failed to parse tenant preferred language config")
		return "", err
	}

	if config.PreferredLanguage == "" {
		return "id", nil
	}

	return config.PreferredLanguage, nil
}
