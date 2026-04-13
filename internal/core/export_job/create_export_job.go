package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *exportJobCore) CreateExportJob(ctx context.Context, data coreentity.ExportJobCreate) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:create_export_job:CreateExportJob")
	defer span.End()

	if !data.UserCtx.HasRole("owner") && !data.UserCtx.HasRole("admin") {
		return nil, errmsg.NewCustomErrors(403).SetMessage("You are not allowed to export attendance reports")
	}

	if data.ResourceType == "" || data.ProcessorKey == "" {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Export job resource metadata is required")
	}

	if c.bus == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage("Export job message bus is not configured")
	}

	item, err := c.repo.CreateExportJob(ctx, data)
	if err != nil {
		return nil, err
	}

	err = c.bus.PublishJSON(ctx, common.MessageSubjectExportJobRequested, coreentity.ExportJobRequestedEvent{
		JobID:    item.ID,
		TenantID: item.TenantID,
	})
	if err != nil {
		errorMessage := err.Error()
		_ = c.repo.UpdateExportJob(ctx, coreentity.ExportJobUpdate{
			TenantID:     item.TenantID,
			ID:           item.ID,
			Status:       coreentity.ExportJobStatusFailed,
			ErrorMessage: &errorMessage,
		})
		return nil, err
	}

	return item, nil
}
