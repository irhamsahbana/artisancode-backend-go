package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) InsertTenant(ctx context.Context, tenant coreentity.Tenant) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:InsertTenant")
	defer span.End()

	query := `
		INSERT INTO tenants (name, code)
		VALUES (?, ?)
		RETURNING id
	`

	var tenantID string
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &tenantID, exec.Rebind(query), tenant.Name, tenant.Code)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, tenant).Msg("Failed to insert tenant")
		return "", err
	}
	return tenantID, nil
}
