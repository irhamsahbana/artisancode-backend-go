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

func (r *userRepo) GetRoleByName(ctx context.Context, roleName, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:register_owner:GetRoleByName")
	defer span.End()

	type dao struct {
		ID       string  `db:"id"`
		TenantID *string `db:"tenant_id"`
		Name     string  `db:"name"`
	}
	var row dao

	query := `
		SELECT id, tenant_id, name
		FROM roles
		WHERE name = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`
	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &row, exec.Rebind(query), roleName, tenantID)
	if err != nil {
		payload := map[string]string{"roleName": roleName, "tenantID": tenantID}
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Role not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("Role not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get role")
		return nil, err
	}

	return &coreentity.Role{
		ID:       row.ID,
		TenantID: row.TenantID,
		Name:     row.Name,
	}, nil
}
