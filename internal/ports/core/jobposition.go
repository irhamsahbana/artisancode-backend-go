package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type JobPositionCore interface {
	GetJobPositions(ctx context.Context, filter coreentity.JobPositionListFilter) ([]coreentity.JobPosition, int, error)
	GetJobPosition(ctx context.Context, filter coreentity.JobPosition) (*coreentity.JobPosition, error)
	CreateJobPosition(ctx context.Context, data coreentity.JobPosition) (*coreentity.JobPosition, error)
	UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error
	DeleteJobPosition(ctx context.Context, filter coreentity.JobPositionDeleteFilter) error
}
