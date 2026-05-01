package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *exportJobCore) GetExportJob(
	ctx context.Context,
	filter coreentity.ExportJobDetailFilter,
) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:get_export_job:GetExportJob")
	defer span.End()

	if !filter.UserCtx.HasRole("owner") && !filter.UserCtx.HasRole("admin") {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAllowedToAccessExportJobs)
	}

	item, err := c.repo.GetExportJob(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = c.attachDownloadURL(ctx, item, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	return item, nil
}
