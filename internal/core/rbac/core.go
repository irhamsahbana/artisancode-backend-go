package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	repository "codebase-app/internal/ports/secondary/db"
)

type rbacCore struct {
	repo repository.RbacRepository
}

var _ corePorts.RbacCore = &rbacCore{}

type Config struct {
	Repo repository.RbacRepository
}

func NewRbacCore(cfg Config) *rbacCore {
	return &rbacCore{
		repo: cfg.Repo,
	}
}

func (c *rbacCore) GetRoles(ctx context.Context, filter coreentity.RoleListFilter) ([]coreentity.Role, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetRoles")
	defer span.End()

	return c.repo.GetRoles(ctx, filter)
}

func (c *rbacCore) GetRoleByName(ctx context.Context, name, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetRoleByName")
	defer span.End()

	return c.repo.GetRoleByName(ctx, name, tenantID)
}

func (c *rbacCore) CreateRole(ctx context.Context, data coreentity.Role) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:CreateRole")
	defer span.End()

	return c.repo.CreateRole(ctx, data)
}

func (c *rbacCore) UpdateRole(ctx context.Context, roleID, tenantID string, permissionIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:UpdateRole")
	defer span.End()

	return c.repo.UpdateRole(ctx, roleID, tenantID, permissionIDs)
}

func (c *rbacCore) DeleteRole(ctx context.Context, filter coreentity.RoleDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:DeleteRole")
	defer span.End()

	return c.repo.DeleteRole(ctx, filter)
}

func (c *rbacCore) GetUserRoles(ctx context.Context, filter coreentity.UserRoleFilter) ([]coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetUserRoles")
	defer span.End()

	return c.repo.GetUserRoles(ctx, filter)
}

func (c *rbacCore) AssignRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:AssignRole")
	defer span.End()

	return c.repo.AssignRole(ctx, data)
}

func (c *rbacCore) RemoveRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:RemoveRole")
	defer span.End()

	return c.repo.RemoveRole(ctx, data)
}

func (c *rbacCore) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:HasPermission")
	defer span.End()

	return c.repo.HasPermission(ctx, userID, permission)
}

func (c *rbacCore) GetRoleWithPermissions(ctx context.Context, roleID, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetRoleWithPermissions")
	defer span.End()

	return c.repo.GetRoleWithPermissions(ctx, roleID, tenantID)
}

func (c *rbacCore) SetRolePermissions(ctx context.Context, roleID, tenantID string, permissionIDs []string) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:SetRolePermissions")
	defer span.End()

	return c.repo.SetRolePermissions(ctx, roleID, tenantID, permissionIDs)
}

func (c *rbacCore) GetPermissions(
	ctx context.Context,
	filter coreentity.PermissionListFilter,
) ([]coreentity.Permission, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetPermissions")
	defer span.End()

	return c.repo.GetPermissions(ctx, filter)
}

func (c *rbacCore) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetPermissionsByRoleID")
	defer span.End()

	return c.repo.GetPermissionsByRoleID(ctx, roleID)
}

func (c *rbacCore) GetUserPermissions(ctx context.Context, userID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:rbac:core:GetUserPermissions")
	defer span.End()

	return c.repo.GetUserPermissions(ctx, userID)
}
