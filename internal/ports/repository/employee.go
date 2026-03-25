package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type EmployeeRepository interface {
	GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error)
	GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error)
	CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error)
	UpdateEmployee(ctx context.Context, data coreentity.Employee) error
	DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error
}