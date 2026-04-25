package coreentity

import "codebase-app/internal/entity/common"

type InternalClient struct {
	UserCtx common.UserContext

	ID          string
	Name        string
	Code        string
	OwnerNames  string
	OwnerEmails string
	CreatedAt   string
	UpdatedAt   string
}

type InternalClientListFilter struct {
	Q        string
	Owner    string
	Page     int
	Paginate int
}

type InternalClientOwnerPermissions struct {
	ClientID           string
	Available          []Permission
	OwnerPermissionIDs []string
}

type InternalClientOwnerPermissionUpdate struct {
	ClientID      string
	PermissionIDs []string
}
