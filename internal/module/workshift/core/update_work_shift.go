package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateWorkShift")
	defer span.End()

	return c.repo.UpdateWorkShift(ctx, data)
}