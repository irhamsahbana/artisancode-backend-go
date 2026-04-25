package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *companyRepo) UpdateCompany(ctx context.Context, data coreentity.Company) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:UpdateCompany")
	defer span.End()

	configJSON, err := json.Marshal(data.Config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to marshal config")
		return err
	}

	query := `
		UPDATE org_units
		SET name = ?, code = ?, config = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND category = 'company' AND deleted_at IS NULL
	`

	_, err = r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Name, data.Code, configJSON, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update company")
		return err
	}
	return nil
}
