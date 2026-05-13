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
	AssignUser(ctx context.Context, tenantID, employeeID, userID string) error
	DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error
	ExistsByEmployeeNo(ctx context.Context, tenantID, employeeNo, excludeID string) (bool, error)
	CountEmployeesByTenant(ctx context.Context, tenantID string) (int64, error)
}
