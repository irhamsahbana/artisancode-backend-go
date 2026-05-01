package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) copyTemplateRolePermissions(ctx context.Context, tx sqlExecutor, tenantID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:initialize_tenant:copyTemplateRolePermissions",
	)
	defer span.End()

	type templateRole struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	type templatePermission struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	type rolePermission struct {
		RoleID       string `db:"role_id"`
		PermissionID string `db:"permission_id"`
	}

	roles := []templateRole{}
	err := tx.SelectContext(ctx, &roles, `SELECT id, name FROM internal_template_roles WHERE deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).
			Msg("Failed to select template roles for role permissions")
		return err
	}

	permissions := []templatePermission{}
	err = tx.SelectContext(
		ctx,
		&permissions,
		`SELECT id, name FROM internal_template_permissions WHERE deleted_at IS NULL`,
	)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).
			Msg("Failed to select template permissions for role permissions")
		return err
	}

	roleIDMapping := make(map[string]string)
	for _, role := range roles {
		var newRoleID string
		err = tx.GetContext(
			ctx,
			&newRoleID,
			tx.Rebind(`SELECT id FROM roles WHERE tenant_id = ? AND name = ?`),
			tenantID,
			role.Name,
		)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "roleName": role.Name}).
				Msg("Failed to get copied role id")
			return err
		}
		roleIDMapping[role.ID] = newRoleID
	}

	permissionIDMapping := make(map[string]string)
	for _, perm := range permissions {
		var newPermID string
		err = tx.GetContext(
			ctx,
			&newPermID,
			tx.Rebind(`SELECT id FROM permissions WHERE tenant_id = ? AND name = ?`),
			tenantID,
			perm.Name,
		)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "permissionName": perm.Name}).
				Msg("Failed to get copied permission id")
			return err
		}
		permissionIDMapping[perm.ID] = newPermID
	}

	rolePermissions := []rolePermission{}
	err = tx.SelectContext(
		ctx,
		&rolePermissions,
		`SELECT role_id, permission_id FROM internal_template_role_permissions`,
	)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).
			Msg("Failed to select template role permissions")
		return err
	}

	for _, rp := range rolePermissions {
		newRoleID := roleIDMapping[rp.RoleID]
		newPermID := permissionIDMapping[rp.PermissionID]
		if newRoleID == "" || newPermID == "" {
			continue
		}

		_, err = tx.ExecContext(ctx, tx.Rebind(`
			INSERT INTO role_permissions (role_id, permission_id, tenant_id)
			VALUES (?, ?, ?)
			ON CONFLICT (role_id, permission_id) DO NOTHING
		`), newRoleID, newPermID, tenantID)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "roleID": newRoleID, "permissionID": newPermID}).
				Msg("Failed to insert role permission")
			return err
		}
	}

	return nil
}
