package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) copyTemplateRoles(ctx context.Context, tx sqlExecutor, tenantID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:copyTemplateRoles",
	)
	defer span.End()

	payload := map[string]string{"tenantID": tenantID}

	type templateRole struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	roles := []templateRole{}
	err := tx.SelectContext(ctx, &roles, `SELECT id, name FROM internal_template_roles WHERE deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to select template roles")
		return err
	}

	for _, role := range roles {
		var newRoleID string
		err = tx.GetContext(ctx, &newRoleID, tx.Rebind(`
			INSERT INTO roles (tenant_id, name)
			VALUES (?, ?)
			ON CONFLICT (tenant_id, name) DO UPDATE SET
				name = EXCLUDED.name,
				updated_at = CURRENT_TIMESTAMP,
				deleted_at = NULL
			RETURNING id
		`), tenantID, role.Name)
		if err != nil {
			payload["roleName"] = role.Name
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to copy role")
			return err
		}
	}

	return nil
}
