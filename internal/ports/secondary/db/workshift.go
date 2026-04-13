package repository

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type WorkShiftRepository interface {
	GetWorkShifts(ctx context.Context, filter coreentity.WorkShiftListFilter) ([]coreentity.WorkShift, int, error)
	GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error)
	CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error)
	UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error
	DeleteWorkShift(ctx context.Context, filter coreentity.WorkShiftDeleteFilter) error
}