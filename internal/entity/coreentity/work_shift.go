package coreentity

import "codebase-app/internal/entity/common"

type WorkShift struct {
	UserCtx common.UserContext

	ID                 string
	TenantID           string
	Name               string
	StartTime          string
	EndTime            string
	GracePeriodMinutes int
}

type WorkShiftListFilter struct {
	TenantID string
	Q        string
	Page     int
	Paginate int
}

type WorkShiftDeleteFilter struct {
	TenantID string
	ID       string
}
