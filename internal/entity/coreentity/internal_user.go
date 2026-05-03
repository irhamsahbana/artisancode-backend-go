package coreentity

import "codebase-app/internal/entity/common"

const (
	InternalUserRoleSuperAdmin = "super_admin"
	InternalUserRoleOperator   = "operator"
	InternalUserRoleFinance    = "finance"
	InternalUserRoleOperations = "operations"

	InternalUserStatusInvited  = "invited"
	InternalUserStatusActive   = "active"
	InternalUserStatusInactive = "inactive"
)

type InternalUser struct {
	UserCtx common.UserContext

	ID           string
	FullName     string
	Email        string
	Password     string
	RefreshToken string
	RoleCode     string
	Status       string
	LastLoginAt  *string
	CreatedAt    string
	UpdatedAt    string
}

type InternalUserListFilter struct {
	Q        string
	RoleCode string
	Status   string
	Page     int
	Paginate int
}

type InternalUserFilter struct {
	ID string
}

type InternalUserDeleteFilter struct {
	ID string
}
