package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) copyTemplatePermissions(ctx context.Context, tx sqlExecutor, tenantID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:copyTemplatePermissions",
	)
	defer span.End()

	type templatePermission struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}

	permissions := []templatePermission{}
	err := tx.SelectContext(
		ctx,
		&permissions,
		`SELECT id, name, COALESCE(description, '') AS description FROM internal_template_permissions WHERE deleted_at IS NULL`,
	)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).
			Msg("Failed to select template permissions")
		return err
	}

	for _, perm := range permissions {
		var newPermID string
		err = tx.GetContext(ctx, &newPermID, tx.Rebind(`
			INSERT INTO permissions (tenant_id, name, description)
			VALUES (?, ?, ?)
			ON CONFLICT (tenant_id, name) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				updated_at = CURRENT_TIMESTAMP,
				deleted_at = NULL
			RETURNING id
		`), tenantID, perm.Name, perm.Description)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "permissionName": perm.Name}).
				Msg("Failed to copy permission")
			return err
		}
	}

	return nil
}
