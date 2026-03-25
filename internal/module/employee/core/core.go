package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
	"context"
)

var _ corePorts.EmployeeCore = &employeeCore{}

type employeeCore struct {
	repo portsRepo.EmployeeRepository
}

type EmployeeCoreConfig struct {
	Repo portsRepo.EmployeeRepository
}

func NewEmployeeCore(cfg EmployeeCoreConfig) *employeeCore {
	return &employeeCore{repo: cfg.Repo}
}

func (c *employeeCore) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetEmployees")
	defer span.End()

	return c.repo.GetEmployees(ctx, filter)
}

func (c *employeeCore) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetEmployee")
	defer span.End()

	return c.repo.GetEmployee(ctx, filter)
}

func (c *employeeCore) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateEmployee")
	defer span.End()

	return c.repo.CreateEmployee(ctx, data)
}

func (c *employeeCore) UpdateEmployee(ctx context.Context, data coreentity.Employee) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateEmployee")
	defer span.End()

	return c.repo.UpdateEmployee(ctx, data)
}

func (c *employeeCore) DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteEmployee")
	defer span.End()

	return c.repo.DeleteEmployee(ctx, filter)
}