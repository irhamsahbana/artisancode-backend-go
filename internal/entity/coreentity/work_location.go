package coreentity

import "codebase-app/internal/entity/common"

type WorkLocation struct {
	UserCtx common.UserContext

	ID           string
	TenantID     string
	OrgUnitID    *string
	OrgUnitName  *string
	Name         string
	Address      *string
	Latitude     *float64
	Longitude    *float64
	RadiusMeters *int
}

type WorkLocationListFilter struct {
	TenantID  string
	OrgUnitID *string
	Q         string
	Page      int
	Paginate  int
}

type WorkLocationDeleteFilter struct {
	TenantID string
	ID       string
}
