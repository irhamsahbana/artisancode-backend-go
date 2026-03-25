package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
	"context"
)

type jobPositionCore struct {
	repo portsRepo.JobPositionRepository
}

type JobPositionCoreConfig struct {
	Repo portsRepo.JobPositionRepository
}

var _ corePorts.JobPositionCore = &jobPositionCore{}

func NewJobPositionCore(cfg JobPositionCoreConfig) *jobPositionCore {
	return &jobPositionCore{repo: cfg.Repo}
}

func (c *jobPositionCore) GetJobPositions(ctx context.Context, filter coreentity.JobPositionListFilter) ([]coreentity.JobPosition, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetJobPositions")
	defer span.End()

	return c.repo.GetJobPositions(ctx, filter)
}

func (c *jobPositionCore) GetJobPosition(ctx context.Context, filter coreentity.JobPosition) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetJobPosition")
	defer span.End()

	return c.repo.GetJobPosition(ctx, filter)
}

func (c *jobPositionCore) CreateJobPosition(ctx context.Context, data coreentity.JobPosition) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateJobPosition")
	defer span.End()

	return c.repo.CreateJobPosition(ctx, data)
}

func (c *jobPositionCore) UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateJobPosition")
	defer span.End()

	return c.repo.UpdateJobPosition(ctx, data)
}

func (c *jobPositionCore) DeleteJobPosition(ctx context.Context, filter coreentity.JobPositionDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteJobPosition")
	defer span.End()

	return c.repo.DeleteJobPosition(ctx, filter)
}
