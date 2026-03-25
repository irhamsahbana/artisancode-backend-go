package coreentity

import "codebase-app/internal/entity/common"

type JobPosition struct {
	UserCtx common.UserContext

	ID       string
	TenantID string
	Name     string
	Grade    *string
}

type JobPositionListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type JobPositionDeleteFilter struct {
	TenantID string
	ID       string
}
