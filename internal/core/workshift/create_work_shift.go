package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:workshift:create_work_shift:CreateWorkShift")
	defer span.End()

	return c.repo.CreateWorkShift(ctx, data)
}
