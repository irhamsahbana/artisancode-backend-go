package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *companyRepo) GetCompany(ctx context.Context, filter coreentity.Company) (*coreentity.Company, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:company:repo:GetCompany")
	defer span.End()

	type dao struct {
		ID        string  `db:"id"`
		TenantID  string  `db:"tenant_id"`
		Name      string  `db:"name"`
		Code      string  `db:"code"`
		Config    []byte  `db:"config"`
		CreatedAt string  `db:"created_at"`
		UpdatedAt *string `db:"updated_at"`
	}

	query := `
		SELECT id, tenant_id, name, code, config, created_at, updated_at
		FROM org_units
		WHERE id = ? AND tenant_id = ? AND category = 'company' AND deleted_at IS NULL
	`

	var data dao
	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageCompanyNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCompanyNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get company")
		return nil, err
	}

	var config coreentity.CompanyConfig
	if len(data.Config) > 0 {
		err = json.Unmarshal(data.Config, &config)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data.ID).Msg("Failed to unmarshal config")
			return nil, err
		}
	}

	updatedAt := ""
	if data.UpdatedAt != nil {
		updatedAt = *data.UpdatedAt
	}

	result := &coreentity.Company{
		ID:        data.ID,
		TenantID:  data.TenantID,
		Name:      data.Name,
		Code:      data.Code,
		Config:    config,
		CreatedAt: data.CreatedAt,
		UpdatedAt: updatedAt,
	}
	return result, nil
}
