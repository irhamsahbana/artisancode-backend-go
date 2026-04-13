package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *meCore) GetMyEmployee(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:me:get_my_employee:GetMyEmployee")
	defer span.End()

	return c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserID)
}
