package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) InitializeTenant(ctx context.Context, tenantID string, companyName string, preferredLanguage string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:InitializeTenant")
	defer span.End()

	tx := r.executor(ctx)

	companyID, err := r.insertDefaultCompany(ctx, tx, tenantID, companyName, preferredLanguage)
	if err != nil {
		return "", err
	}

	err = r.copyRoles(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	err = r.copyPermissions(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	err = r.copyRolePermissions(ctx, tx, tenantID)
	if err != nil {
		return "", err
	}

	return companyID, nil
}

func (r *userRepo) insertDefaultCompany(ctx context.Context, tx interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(string) string
}, tenantID string, companyName string, preferredLanguage string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:insertDefaultCompany")
	defer span.End()

	config := coreentity.DefaultCompanyConfig()
	if preferredLanguage != "" {
		config.PreferredLanguage = preferredLanguage
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "companyName": companyName}).Msg("Failed to marshal company config")
		return "", err
	}

	query := `
		INSERT INTO org_units (tenant_id, name, code, category, config)
		VALUES (?, ?, ?, 'company', ?)
		RETURNING id
	`
	var companyID string
	err = tx.GetContext(ctx, &companyID, tx.Rebind(query), tenantID, companyName, companyName, configJSON)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "companyName": companyName}).Msg("Failed to insert default company")
		return "", err
	}
	return companyID, nil
}

func (r *userRepo) copyRoles(ctx context.Context, tx interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(string) string
}, tenantID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:copyRoles")
	defer span.End()

	payload := map[string]string{"tenantID": tenantID}

	type templateRole struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	roles := []templateRole{}
	err := tx.SelectContext(ctx, &roles, `SELECT id, name FROM roles WHERE tenant_id IS NULL AND deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to select template roles")
		return err
	}

	roleIDMapping := make(map[string]string)

	for _, role := range roles {
		var newRoleID string
		err = tx.GetContext(ctx, &newRoleID, tx.Rebind(`
			INSERT INTO roles (tenant_id, name)
			VALUES (?, ?)
			ON CONFLICT (tenant_id, name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`), tenantID, role.Name)
		if err != nil {
			payload["roleName"] = role.Name
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to copy role")
			return err
		}
		roleIDMapping[role.ID] = newRoleID
	}

	return nil
}

func (r *userRepo) copyPermissions(ctx context.Context, tx interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(string) string
}, tenantID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:copyPermissions")
	defer span.End()

	type templatePermission struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}

	permissions := []templatePermission{}
	err := tx.SelectContext(ctx, &permissions, `SELECT id, name, description FROM permissions WHERE tenant_id IS NULL AND deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).Msg("Failed to select template permissions")
		return err
	}

	permissionIDMapping := make(map[string]string)

	for _, perm := range permissions {
		var newPermID string
		err = tx.GetContext(ctx, &newPermID, tx.Rebind(`
			INSERT INTO permissions (tenant_id, name, description)
			VALUES (?, ?, ?)
			ON CONFLICT (tenant_id, name) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description
			RETURNING id
		`), tenantID, perm.Name, perm.Description)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "permissionName": perm.Name}).Msg("Failed to copy permission")
			return err
		}
		permissionIDMapping[perm.ID] = newPermID
	}

	return nil
}

func (r *userRepo) copyRolePermissions(ctx context.Context, tx interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(string) string
}, tenantID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:initialize_tenant:copyRolePermissions")
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
	err := tx.SelectContext(ctx, &roles, `SELECT id, name FROM roles WHERE tenant_id IS NULL AND deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).Msg("Failed to select template roles for role permissions")
		return err
	}

	permissions := []templatePermission{}
	err = tx.SelectContext(ctx, &permissions, `SELECT id, name FROM permissions WHERE tenant_id IS NULL AND deleted_at IS NULL`)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).Msg("Failed to select template permissions for role permissions")
		return err
	}

	roleIDMapping := make(map[string]string)
	for _, role := range roles {
		var newRoleID string
		err = tx.GetContext(ctx, &newRoleID, tx.Rebind(`SELECT id FROM roles WHERE tenant_id = ? AND name = ?`), tenantID, role.Name)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "roleName": role.Name}).Msg("Failed to get copied role id")
			return err
		}
		roleIDMapping[role.ID] = newRoleID
	}

	permissionIDMapping := make(map[string]string)
	for _, perm := range permissions {
		var newPermID string
		err = tx.GetContext(ctx, &newPermID, tx.Rebind(`SELECT id FROM permissions WHERE tenant_id = ? AND name = ?`), tenantID, perm.Name)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "permissionName": perm.Name}).Msg("Failed to get copied permission id")
			return err
		}
		permissionIDMapping[perm.ID] = newPermID
	}

	rolePermissions := []rolePermission{}
	err = tx.SelectContext(ctx, &rolePermissions, tx.Rebind(`SELECT role_id, permission_id FROM role_permissions WHERE tenant_id IS NULL`))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID}).Msg("Failed to select template role permissions")
		return err
	}

	for _, rp := range rolePermissions {
		newRoleID := roleIDMapping[rp.RoleID]
		newPermID := permissionIDMapping[rp.PermissionID]
		if newRoleID != "" && newPermID != "" {
			_, err = tx.ExecContext(ctx, tx.Rebind(`
				INSERT INTO role_permissions (role_id, permission_id, tenant_id)
				VALUES (?, ?, ?)
				ON CONFLICT (role_id, permission_id) DO NOTHING
			`), newRoleID, newPermID, tenantID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"tenantID": tenantID, "roleID": newRoleID, "permissionID": newPermID}).Msg("Failed to insert role permission")
				return err
			}
		}
	}

	return nil
}
