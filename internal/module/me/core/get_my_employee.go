package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *meCore) GetMyEmployee(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetMyEmployee")
	defer span.End()

	return c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserID)
}
