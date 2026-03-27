package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workLocationCore) DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteWorkLocation")
	defer span.End()

	return c.repo.DeleteWorkLocation(ctx, filter)
}