package seeds

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type seedPermission struct {
	Name        string
	Description string
}

func (s *Seed) rbacSeed() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	permissions := []seedPermission{
		{Name: "master.read", Description: "Read master data"},
		{Name: "master.write", Description: "Create/update/delete master data"},
		{Name: "employee.read", Description: "Read employee data"},
		{Name: "employee.write", Description: "Create/update/delete employee data"},
		{Name: "attendance.read", Description: "Read attendance data"},
		{Name: "attendance.write", Description: "Manage attendance data"},
		{Name: "user.read", Description: "Read user data"},
		{Name: "user.write", Description: "Create/update/delete user data"},
		{Name: "rbac.manage", Description: "Manage roles and permissions"},
		{Name: "report.read", Description: "Read reports"},
	}
	var defaultTenantID *string

	rolePermissions := map[string][]string{
		"owner":    {"master.read", "master.write", "employee.read", "employee.write", "attendance.read", "attendance.write", "user.read", "user.write", "rbac.manage", "report.read"},
		"admin":    {"master.read", "master.write", "employee.read", "employee.write", "attendance.read", "attendance.write", "user.read", "user.write", "report.read"},
		"manager":  {"master.read", "employee.read", "attendance.read", "report.read"},
		"employee": {"master.read", "attendance.read"},
	}

	permissionIDs := make(map[string]string, len(permissions))
	for _, perm := range permissions {
		id, err := upsertPermission(tx, defaultTenantID, perm)
		if err != nil {
			return
		}
		permissionIDs[perm.Name] = id
	}

	for roleName, perms := range rolePermissions {
		roleID, err := upsertRole(tx, defaultTenantID, roleName)
		if err != nil {
			return
		}
		for _, permName := range perms {
			permID, ok := permissionIDs[permName]
			if !ok {
				continue
			}
			if err := insertRolePermission(tx, roleID, permID); err != nil {
				return
			}
		}
	}

	log.Info().Msg("rbac seeded successfully!")
}

func upsertRole(tx *sqlx.Tx, tenantID *string, name string) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM roles
		WHERE name = ? AND tenant_id IS NOT DISTINCT FROM ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), name, tenantID); err == nil {
		return id, nil
	}

	query := `
		INSERT INTO roles (tenant_id, name)
		VALUES (?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(query), tenantID, name); err != nil {
		log.Error().Err(err).Any("name", name).Msg("failed to upsert role")
		return "", err
	}
	return id, nil
}

func upsertPermission(tx *sqlx.Tx, tenantID *string, perm seedPermission) (string, error) {
	var id string
	selectQuery := `
		SELECT id
		FROM permissions
		WHERE name = ? AND tenant_id IS NOT DISTINCT FROM ? AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`
	if err := tx.Get(&id, tx.Rebind(selectQuery), perm.Name, tenantID); err == nil {
		return id, nil
	}

	query := `
		INSERT INTO permissions (tenant_id, name, description)
		VALUES (?, ?, ?)
		RETURNING id
	`
	if err := tx.Get(&id, tx.Rebind(query), tenantID, perm.Name, perm.Description); err != nil {
		log.Error().Err(err).Any("name", perm.Name).Msg("failed to upsert permission")
		return "", err
	}
	return id, nil
}

func insertRolePermission(tx *sqlx.Tx, roleID, permissionID string) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`
	if _, err := tx.Exec(tx.Rebind(query), roleID, permissionID); err != nil {
		log.Error().Err(err).Msg("failed to insert role permission")
		return err
	}
	return nil
}
