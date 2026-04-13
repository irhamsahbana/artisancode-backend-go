package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:workshift:update_work_shift:UpdateWorkShift")
	defer span.End()

	return c.repo.UpdateWorkShift(ctx, data)
}
