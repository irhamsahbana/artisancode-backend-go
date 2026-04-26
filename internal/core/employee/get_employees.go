package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *employeeCore) GetEmployees(
	ctx context.Context,
	filter coreentity.EmployeeListFilter,
) ([]coreentity.Employee, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:employee:get_employees:GetEmployees")
	defer span.End()

	return c.repo.GetEmployees(ctx, filter)
}
