package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *employeeCore) DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:employee:delete_employee:DeleteEmployee")
	defer span.End()

	return c.repo.DeleteEmployee(ctx, filter)
}
