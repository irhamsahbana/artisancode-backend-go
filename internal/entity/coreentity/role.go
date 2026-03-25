package coreentity

import "codebase-app/internal/entity/common"

type Role struct {
	UserCtx common.UserContext

	ID       string
	TenantID *string
	Name     string
}

type RoleListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type RoleDeleteFilter struct {
	TenantID string
	ID       string
}

type UserRole struct {
	UserCtx common.UserContext

	UserID string
	RoleID string
}

type UserRoleFilter struct {
	UserID string
}

type Permission struct {
	ID   string
	Name string
}

type UserPermission struct {
	UserID         string
	PermissionID   string
	PermissionName string
}
