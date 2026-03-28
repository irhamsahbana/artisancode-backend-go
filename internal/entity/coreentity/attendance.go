package coreentity

import "codebase-app/internal/entity/common"

type AttendanceLog struct {
	UserCtx common.UserContext

	ID             string
	TenantID       string
	EmployeeID     string
	EmployeeNo     string
	EmployeeName   string
	AttendanceDate string
	Type           string
	Source         string
	Status         string
	LoggedAt       string
	Latitude       *float64
	Longitude      *float64
	Address        *string
	DeviceID       *string
	DeviceName     *string
	Notes          *string
	CreatedAt      string
	UpdatedAt      *string
}

type AttendanceLogListFilter struct {
	UserCtx       common.UserContext
	TenantID      string
	EmployeeID    *string
	Q             string
	Type          string
	Source        string
	AttendanceDay *string
	DateFrom      *string
	DateTo        *string
	Page          int
	Paginate      int
}

type AttendanceLogDetailFilter struct {
	UserCtx    common.UserContext
	TenantID   string
	ID         string
	EmployeeID *string
}

type AttendanceLogAction struct {
	UserCtx    common.UserContext
	TenantID   string
	Type       string
	LoggedAt   string
	Latitude   *float64
	Longitude  *float64
	Address    *string
	DeviceID   *string
	DeviceName *string
	Notes      *string
}
