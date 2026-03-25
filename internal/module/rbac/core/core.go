package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/internal/ports/repository"
)

type rbacCore struct {
	repo repository.RbacRepository
}

var _ corePorts.RbacCore = &rbacCore{}

func NewRbacCore(repo repository.RbacRepository) *rbacCore {
	return &rbacCore{
		repo: repo,
	}
}

func (c *rbacCore) GetRoles(ctx context.Context, filter coreentity.RoleListFilter) ([]coreentity.Role, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetRoles")
	defer span.End()

	return c.repo.GetRoles(ctx, filter)
}

func (c *rbacCore) GetRoleByName(ctx context.Context, name, tenantID string) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetRoleByName")
	defer span.End()

	return c.repo.GetRoleByName(ctx, name, tenantID)
}

func (c *rbacCore) CreateRole(ctx context.Context, data coreentity.Role) (*coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateRole")
	defer span.End()

	return c.repo.CreateRole(ctx, data)
}

func (c *rbacCore) UpdateRole(ctx context.Context, data coreentity.Role) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateRole")
	defer span.End()

	return c.repo.UpdateRole(ctx, data)
}

func (c *rbacCore) DeleteRole(ctx context.Context, filter coreentity.RoleDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteRole")
	defer span.End()

	return c.repo.DeleteRole(ctx, filter)
}

func (c *rbacCore) GetUserRoles(ctx context.Context, filter coreentity.UserRoleFilter) ([]coreentity.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetUserRoles")
	defer span.End()

	return c.repo.GetUserRoles(ctx, filter)
}

func (c *rbacCore) AssignRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "core.AssignRole")
	defer span.End()

	return c.repo.AssignRole(ctx, data)
}

func (c *rbacCore) RemoveRole(ctx context.Context, data coreentity.UserRole) error {
	ctx, span := tracing.StartSpan(ctx, "core.RemoveRole")
	defer span.End()

	return c.repo.RemoveRole(ctx, data)
}

func (c *rbacCore) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "core.HasPermission")
	defer span.End()

	return c.repo.HasPermission(ctx, userID, permission)
}

func (c *rbacCore) GetUserPermissions(ctx context.Context, userID string) ([]coreentity.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetUserPermissions")
	defer span.End()

	return c.repo.GetUserPermissions(ctx, userID)
}