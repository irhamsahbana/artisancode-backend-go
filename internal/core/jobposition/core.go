package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"
)

type jobPositionCore struct {
	repo portsRepo.JobPositionRepository
}

type Config struct {
	Repo portsRepo.JobPositionRepository
}

var _ corePorts.JobPositionCore = &jobPositionCore{}

func NewJobPositionCore(cfg Config) *jobPositionCore {
	return &jobPositionCore{repo: cfg.Repo}
}

func (c *jobPositionCore) GetJobPositions(
	ctx context.Context,
	filter coreentity.JobPositionListFilter,
) ([]coreentity.JobPosition, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:jobposition:core:GetJobPositions")
	defer span.End()

	return c.repo.GetJobPositions(ctx, filter)
}

func (c *jobPositionCore) GetJobPosition(
	ctx context.Context,
	filter coreentity.JobPosition,
) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:jobposition:core:GetJobPosition")
	defer span.End()

	return c.repo.GetJobPosition(ctx, filter)
}

func (c *jobPositionCore) CreateJobPosition(
	ctx context.Context,
	data coreentity.JobPosition,
) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:jobposition:core:CreateJobPosition")
	defer span.End()

	return c.repo.CreateJobPosition(ctx, data)
}

func (c *jobPositionCore) UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:jobposition:core:UpdateJobPosition")
	defer span.End()

	return c.repo.UpdateJobPosition(ctx, data)
}

func (c *jobPositionCore) DeleteJobPosition(ctx context.Context, filter coreentity.JobPositionDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:jobposition:core:DeleteJobPosition")
	defer span.End()

	return c.repo.DeleteJobPosition(ctx, filter)
}
