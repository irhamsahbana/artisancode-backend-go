package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *workLocationCore) GetWorkLocations(ctx context.Context, filter coreentity.WorkLocationListFilter) ([]coreentity.WorkLocation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetWorkLocations")
	defer span.End()

	return c.repo.GetWorkLocations(ctx, filter)
}

func (c *workLocationCore) GetWorkLocation(ctx context.Context, filter coreentity.WorkLocation) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetWorkLocation")
	defer span.End()

	return c.repo.GetWorkLocation(ctx, filter)
}