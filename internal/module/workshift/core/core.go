package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
	"context"
)

type workShiftCore struct {
	repo portsRepo.WorkShiftRepository
}

type WorkShiftCoreConfig struct {
	Repo portsRepo.WorkShiftRepository
}

var _ corePorts.WorkShiftCore = &workShiftCore{}

func NewWorkShiftCore(cfg WorkShiftCoreConfig) *workShiftCore {
	return &workShiftCore{repo: cfg.Repo}
}

func (c *workShiftCore) GetWorkShifts(ctx context.Context, filter coreentity.WorkShiftListFilter) ([]coreentity.WorkShift, int, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetWorkShifts")
	defer span.End()

	return c.repo.GetWorkShifts(ctx, filter)
}

func (c *workShiftCore) GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetWorkShift")
	defer span.End()

	return c.repo.GetWorkShift(ctx, filter)
}

func (c *workShiftCore) CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CreateWorkShift")
	defer span.End()

	return c.repo.CreateWorkShift(ctx, data)
}

func (c *workShiftCore) UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error {
	ctx, span := tracing.StartSpan(ctx, "core.UpdateWorkShift")
	defer span.End()

	return c.repo.UpdateWorkShift(ctx, data)
}

func (c *workShiftCore) DeleteWorkShift(ctx context.Context, filter coreentity.WorkShiftDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "core.DeleteWorkShift")
	defer span.End()

	return c.repo.DeleteWorkShift(ctx, filter)
}
