package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetWorkShift")
	defer span.End()

	return c.repo.GetWorkShift(ctx, filter)
}