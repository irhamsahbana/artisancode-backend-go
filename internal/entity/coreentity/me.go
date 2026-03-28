package coreentity

import "codebase-app/internal/entity/common"

type Me struct {
	UserCtx common.UserContext

	UserID      string
	UserName    string
	TenantID    string
	TenantName  string
	Roles       []string
	CompanyID   *string
	CompanyName *string
}

type SelfFilter struct {
	UserCtx common.UserContext

	TenantID string
	UserID   string
}

type WorkShiftToday struct {
	UserCtx common.UserContext

	ShiftID        string
	ShiftName      string
	StartTime      string
	EndTime        string
	Timezone       string
	AttendanceDate string
}
