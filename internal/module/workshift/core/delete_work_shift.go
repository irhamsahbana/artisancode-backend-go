package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) DeleteWorkShift(ctx context.Context, filter coreentity.WorkShiftDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteWorkShift")
	defer span.End()

	return c.repo.DeleteWorkShift(ctx, filter)
}