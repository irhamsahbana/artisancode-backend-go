package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workShiftCore) GetWorkShifts(
	ctx context.Context,
	filter coreentity.WorkShiftListFilter,
) ([]coreentity.WorkShift, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:workshift:get_work_shifts:GetWorkShifts")
	defer span.End()

	return c.repo.GetWorkShifts(ctx, filter)
}
