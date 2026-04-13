package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
)

type rbacRepo struct {
	db *sqlx.DB
}

var _ portsRepo.RbacRepository = &rbacRepo{}

type RbacRepositoryConfig struct {
	DB *sqlx.DB
}

func NewRbacRepository(cfg RbacRepositoryConfig) portsRepo.RbacRepository {
	return &rbacRepo{
		db: cfg.DB,
	}
}

func (r *rbacRepo) GetRoles(ctx context.Context, filter coreentity.RoleListFilter) ([]coreentity.Role, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetRoles")
	defer span.End()

	type row struct {
		TotalData int     `db:"total_data"`
		ID        string  `db:"id"`
		TenantID  *string `db:"tenant_id"`
		Name      string  `db:"name"`
	}

	rows := make([]row, 0)
	args := make([]any, 0, 4)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			tenant_id,
			name
		FROM roles
		WHERE deleted_at IS NULL AND tenant_id = ?
	`
	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		return nil, 0, err
	}

	items := make([]coreentity.Role, 0, len(rows))
	total := 0
	for _, item := range rows {
		total = item.TotalData
		items = append(items, coreentity.Role{
			ID:       item.ID,
			TenantID: item.TenantID,
			Name:     item.Name,
		})
	}
	return items, total, nil
}

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
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Role not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *rbacRepo) CreateRole(ctx context.Context, data coreentity.Role) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:CreateRole")
	defer span.End()

	query := `
		INSERT INTO roles (tenant_id, name)
		VALUES (?, ?)
		RETURNING id, tenant_id, name
	`

	var item coreentity.Role
	if err := r.db.GetContext(ctx, &item, r.db.Rebind(query), data.TenantID, data.Name); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *rbacRepo) UpdateRole(ctx context.Context, data coreentity.Role) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:UpdateRole")
	defer span.End()

	query := `
		UPDATE roles
		SET name = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.Name, data.ID, data.TenantID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("Role not found")
	}
	return nil
}

func (r *rbacRepo) DeleteRole(ctx context.Context, filter coreentity.RoleDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:DeleteRole")
	defer span.End()

	query := `
		UPDATE roles
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("Role not found")
	}
	return nil
}

func (r *rbacRepo) GetUserRoles(ctx context.Context, filter coreentity.UserRoleFilter) ([]coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetUserRoles")
	defer span.End()

	query := `
		SELECT r.id, r.tenant_id, r.name
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL
		ORDER BY r.name ASC
	`

	rows := make([]coreentity.Role, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), filter.UserID); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *rbacRepo) AssignRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:AssignRole")
	defer span.End()

	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.UserID, data.RoleID)
	return err
}

func (r *rbacRepo) RemoveRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:RemoveRole")
	defer span.End()

	query := `
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.UserID, data.RoleID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("User role not found")
	}
	return nil
}

func (r *rbacRepo) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:HasPermission")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM role_permissions rp
			INNER JOIN roles r ON r.id = rp.role_id
			INNER JOIN permissions p ON p.id = rp.permission_id
			INNER JOIN user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = ? AND p.name = ? AND r.deleted_at IS NULL
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), userID, permission); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *rbacRepo) GetUserPermissions(ctx context.Context, userID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetUserPermissions")
	defer span.End()

	query := `
		SELECT DISTINCT p.id, p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		INNER JOIN roles r ON r.id = rp.role_id
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL
		ORDER BY p.name ASC
	`

	rows := make([]coreentity.Permission, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), userID); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *rbacRepo) GetPermissions(ctx context.Context, filter coreentity.PermissionListFilter) ([]coreentity.Permission, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetPermissions")
	defer span.End()

	type row struct {
		TotalData int `db:"total_data"`
		coreentity.Permission
	}

	rows := make([]row, 0)
	args := make([]any, 0, 4)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			name,
			COALESCE(description, '') AS description
		FROM permissions
		WHERE deleted_at IS NULL
	`
	if filter.TenantID != "" {
		query += ` AND tenant_id = ?`
		args = append(args, filter.TenantID)
	}

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		return nil, 0, err
	}

	items := make([]coreentity.Permission, 0, len(rows))
	total := 0
	for _, item := range rows {
		total = item.TotalData
		items = append(items, item.Permission)
	}
	return items, total, nil
}

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
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Role not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *rbacRepo) SetRolePermissions(ctx context.Context, roleID, tenantID string, permissionIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:SetRolePermissions")
	defer span.End()

	// Delete existing permissions
	deleteQuery := `DELETE FROM role_permissions WHERE role_id = ?`
	if _, err := r.db.ExecContext(ctx, r.db.Rebind(deleteQuery), roleID); err != nil {
		return err
	}

	// Insert new permissions
	if len(permissionIDs) > 0 {
		insertQuery := `INSERT INTO role_permissions (role_id, permission_id, tenant_id) VALUES (?, ?, ?)`
		for _, permID := range permissionIDs {
			if _, err := r.db.ExecContext(ctx, r.db.Rebind(insertQuery), roleID, permID, tenantID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *rbacRepo) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:GetPermissionsByRoleID")
	defer span.End()

	query := `
		SELECT p.id, p.name, COALESCE(p.description, '') AS description
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ? AND p.deleted_at IS NULL
		ORDER BY p.name ASC
	`

	rows := make([]coreentity.Permission, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), roleID); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *rbacRepo) AssignPermission(ctx context.Context, userID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:AssignPermission")
	defer span.End()

	query := `
		INSERT INTO user_permissions (user_id, permission_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, permission_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), userID, permissionID)
	return err
}

func (r *rbacRepo) RemovePermission(ctx context.Context, userID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:rbac:repo:RemovePermission")
	defer span.End()

	query := `
		DELETE FROM user_permissions
		WHERE user_id = ? AND permission_id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), userID, permissionID)
	return err
}
