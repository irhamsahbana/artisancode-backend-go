package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
	"context"
)

type workLocationCore struct {
	repo portsRepo.WorkLocationRepository
}

type WorkLocationCoreConfig struct {
	Repo portsRepo.WorkLocationRepository
}

var _ corePorts.WorkLocationCore = &workLocationCore{}

func NewWorkLocationCore(cfg WorkLocationCoreConfig) *workLocationCore {
	return &workLocationCore{repo: cfg.Repo}
}

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

func (c *workLocationCore) CreateWorkLocation(ctx context.Context, data coreentity.WorkLocation) (*coreentity.WorkLocation, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateWorkLocation")
	defer span.End()

	return c.repo.CreateWorkLocation(ctx, data)
}

func (c *workLocationCore) UpdateWorkLocation(ctx context.Context, data coreentity.WorkLocation) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateWorkLocation")
	defer span.End()

	return c.repo.UpdateWorkLocation(ctx, data)
}

func (c *workLocationCore) DeleteWorkLocation(ctx context.Context, filter coreentity.WorkLocationDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteWorkLocation")
	defer span.End()

	return c.repo.DeleteWorkLocation(ctx, filter)
}
