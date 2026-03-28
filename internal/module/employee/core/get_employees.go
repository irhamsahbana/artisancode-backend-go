package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *employeeCore) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetEmployees")
	defer span.End()

	return c.repo.GetEmployees(ctx, filter)
}