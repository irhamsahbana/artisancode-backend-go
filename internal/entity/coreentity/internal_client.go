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
	Page     int
	Paginate int
}
