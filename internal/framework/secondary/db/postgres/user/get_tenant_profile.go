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

func (r *userRepo) GetTenantProfile(ctx context.Context, tenantID string) (*coreentity.TenantProfile, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:get_tenant_profile:GetTenantProfile",
	)
	defer span.End()

	query := `
		SELECT
			id,
			name,
			code
		FROM tenants
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var row struct {
		ID   string `db:"id"`
		Name string `db:"name"`
		Code string `db:"code"`
	}

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &row, exec.Rebind(query), tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).
				Warn().
				Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID}).
				Msg("Tenant profile not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageTenantNotFound)
		}
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenant_id": tenantID}).
			Msg("Failed to get tenant profile")
		return nil, err
	}

	return &coreentity.TenantProfile{
		ID:                  row.ID,
		Name:                row.Name,
		Code:                row.Code,
		CanChangeTenantCode: false,
	}, nil
}
