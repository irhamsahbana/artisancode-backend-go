package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type ExportJobRepository interface {
	CreateExportJob(ctx context.Context, data coreentity.ExportJobCreate) (*coreentity.ExportJob, error)
	GetExportJobs(ctx context.Context, filter coreentity.ExportJobListFilter) ([]coreentity.ExportJob, int, error)
	GetExportJob(ctx context.Context, filter coreentity.ExportJobDetailFilter) (*coreentity.ExportJob, error)
	ClaimPendingExportJob(ctx context.Context) (*coreentity.ExportJob, error)
	UpdateExportJob(ctx context.Context, data coreentity.ExportJobUpdate) error
}
