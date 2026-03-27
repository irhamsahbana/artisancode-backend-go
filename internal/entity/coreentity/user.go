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
