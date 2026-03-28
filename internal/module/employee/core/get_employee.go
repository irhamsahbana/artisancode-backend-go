package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *employeeCore) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetEmployee")
	defer span.End()

	return c.repo.GetEmployee(ctx, filter)
}