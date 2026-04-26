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

func (r *rbacRepo) GetRoleByName(ctx context.Context, name, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetRoleByName")
	defer span.End()

	query := `
		SELECT id, tenant_id, name
		FROM roles
		WHERE name = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var item coreentity.Role
	if err := r.db.GetContext(ctx, &item, r.db.Rebind(query), name, tenantID); err != nil {
		payload := map[string]string{"name": name, "tenant_id": tenantID}
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Role not found by name")
			return nil, errmsg.NewCustomErrors(404).SetMessage("Role not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get role by name")
		return nil, err
	}
	return &item, nil
}
