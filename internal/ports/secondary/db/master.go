package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type MasterRepository interface {
	GetOrgUnits(ctx context.Context, filter coreentity.OrgUnitListFilter) ([]coreentity.OrgUnit, int, error)
	GetOrgUnit(ctx context.Context, filter coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error)
	UpdateOrgUnit(ctx context.Context, data coreentity.OrgUnit) error
	DeleteOrgUnit(ctx context.Context, filter coreentity.OrgUnitDeleteFilter) error

	GetJobPositions(ctx context.Context, filter coreentity.JobPositionListFilter) ([]coreentity.JobPosition, int, error)
	GetJobPosition(ctx context.Context, filter coreentity.JobPosition) (*coreentity.JobPosition, error)
	CreateJobPosition(ctx context.Context, data coreentity.JobPosition) (*coreentity.JobPosition, error)
	UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error
	DeleteJobPosition(ctx context.Context, filter coreentity.JobPositionDeleteFilter) error

	GetWorkLocations(ctx context.Context, filter coreentity.WorkLocationListFilter) ([]coreentity.WorkLocation, int, error)
	GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error)
	UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error
	DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error

	GetWorkShifts(ctx context.Context, filter coreentity.WorkShiftListFilter) ([]coreentity.WorkShift, int, error)
	GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error)
	CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error)
	UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error
	DeleteWorkShift(ctx context.Context, filter coreentity.WorkShiftDeleteFilter) error

	GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error)
	GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error)
	CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error)
	UpdateEmployee(ctx context.Context, data coreentity.Employee) error
	DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error
}
