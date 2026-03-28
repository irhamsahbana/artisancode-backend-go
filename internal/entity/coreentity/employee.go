package coreentity

import "codebase-app/internal/entity/common"

type Employee struct {
	UserCtx common.UserContext

	ID            string
	TenantID      string
	EmployeeNo    string
	FullName      string
	Email         string
	UserID        *string
	OrgUnitID     *string
	JobPositionID *string
	LocationID    *string
	ShiftID       *string
	Status        string
	JoinDate      *string
}

type EmployeeListFilter struct {
	TenantID  string
	Q         string
	Status    string
	OrgUnitID *string
	Page      int
	Paginate  int
}

type EmployeeDeleteFilter struct {
	TenantID string
	ID       string
}