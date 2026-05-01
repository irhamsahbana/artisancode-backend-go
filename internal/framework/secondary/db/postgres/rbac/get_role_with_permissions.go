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

func (r *rbacRepo) GetRoleWithPermissions(ctx context.Context, roleID, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetRoleWithPermissions")
	defer span.End()

	query := `
		SELECT id, tenant_id, name
		FROM roles
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var item coreentity.Role
	if err := r.db.GetContext(ctx, &item, r.db.Rebind(query), roleID, tenantID); err != nil {
		payload := map[string]string{"role_id": roleID, "tenant_id": tenantID}
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Role not found with permissions")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageRoleNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get role with permissions")
		return nil, err
	}
	return &item, nil
}
