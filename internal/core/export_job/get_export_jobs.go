package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *exportJobCore) GetExportJobs(
	ctx context.Context,
	filter coreentity.ExportJobListFilter,
) ([]coreentity.ExportJob, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:get_export_jobs:GetExportJobs")
	defer span.End()

	if !filter.UserCtx.HasRole("owner") && !filter.UserCtx.HasRole("admin") {
		return nil, 0, errmsg.NewCustomErrors(403).SetMessage("You are not allowed to access export jobs")
	}

	items, total, err := c.repo.GetExportJobs(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for index := range items {
		err = c.attachDownloadURL(ctx, &items[index], time.Now().UTC())
		if err != nil {
			return nil, 0, err
		}
	}

	return items, total, nil
}
