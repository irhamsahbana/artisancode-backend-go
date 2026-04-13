package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type ExportJobCore interface {
	CreateExportJob(ctx context.Context, data coreentity.ExportJobCreate) (*coreentity.ExportJob, error)
	GetExportJobs(ctx context.Context, filter coreentity.ExportJobListFilter) ([]coreentity.ExportJob, int, error)
	GetExportJob(ctx context.Context, filter coreentity.ExportJobDetailFilter) (*coreentity.ExportJob, error)
	ProcessExportJob(ctx context.Context, filter coreentity.ExportJobDetailFilter) error
	ProcessPendingExportJobs(ctx context.Context, limit int) (*coreentity.ExportJobProcessResult, error)
}
