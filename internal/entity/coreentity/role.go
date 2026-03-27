package coreentity

import "codebase-app/internal/entity/common"

type Role struct {
	UserCtx common.UserContext `db:"-"`

	ID       string  `db:"id"`
	TenantID *string `db:"tenant_id"`
	Name     string  `db:"name"`
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
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type PermissionListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type UserPermission struct {
	UserID         string
	PermissionID   string
	PermissionName string
}
