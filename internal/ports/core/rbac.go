package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type RbacCore interface {
	GetRoles(ctx context.Context, filter coreentity.RoleListFilter) ([]coreentity.Role, int, error)
	GetRoleByName(ctx context.Context, name, tenantID string) (*coreentity.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID, tenantID string) (*coreentity.Role, error)
	SetRolePermissions(ctx context.Context, roleID, tenantID string, permissionIDs []string) error
	CreateRole(ctx context.Context, data coreentity.Role) (*coreentity.Role, error)
	UpdateRole(ctx context.Context, data coreentity.Role) error
	DeleteRole(ctx context.Context, filter coreentity.RoleDeleteFilter) error

	GetUserRoles(ctx context.Context, filter coreentity.UserRoleFilter) ([]coreentity.Role, error)
	AssignRole(ctx context.Context, data coreentity.UserRole) error
	RemoveRole(ctx context.Context, data coreentity.UserRole) error

	GetPermissions(ctx context.Context, filter coreentity.PermissionListFilter) ([]coreentity.Permission, int, error)
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]coreentity.Permission, error)
	HasPermission(ctx context.Context, userID, permission string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]coreentity.Permission, error)
}
