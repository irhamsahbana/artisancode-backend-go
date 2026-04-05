package coreentity

import "codebase-app/internal/entity/common"

type User struct {
	UserCtx common.UserContext

	ID           string
	RoleIDs      []string
	RoleNames    []string
	Name         string
	UserName     string
	Email        string
	Password     string
	RefreshToken string
	TenantID     string
	TenantCode   string
	TenantName   string
	CompanyID    *string
	CompanyName  *string
}

type UserListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type UserDeleteFilter struct {
	UserCtx  common.UserContext
	TenantID string
	ID       string
}
