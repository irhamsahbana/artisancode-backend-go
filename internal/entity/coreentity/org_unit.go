package coreentity

import "codebase-app/internal/entity/common"

type OrgUnit struct {
	UserCtx common.UserContext

	ID       string
	TenantID string
	Name     string
	ParentID *string
	Category string
}

type OrgUnitListFilter struct {
	TenantID string
	Q        string
	Category string
	Page     int
	Paginate int
}

type OrgUnitDeleteFilter struct {
	TenantID string
	ID       string
}