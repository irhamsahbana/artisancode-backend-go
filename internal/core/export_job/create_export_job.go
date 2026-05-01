package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *exportJobCore) CreateExportJob(
	ctx context.Context,
	data coreentity.ExportJobCreate,
) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:create_export_job:CreateExportJob")
	defer span.End()

	if !data.UserCtx.HasRole("owner") && !data.UserCtx.HasRole("admin") {
		return nil, errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAllowedToExportAttendanceReports)
	}

	if data.ResourceType == "" || data.ProcessorKey == "" {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageExportJobResourceMetadataIsRequired)
	}

	if c.bus == nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage(errmsg.MessageExportJobMessageBusIsNotConfigured)
	}

	var item *coreentity.ExportJob
	err := c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		createdItem, txErr := c.repo.CreateExportJob(txCtx, data)
		if txErr != nil {
			return txErr
		}
		item = createdItem

		return c.bus.PublishJSON(txCtx, common.MessageTopicExportJobRequested, coreentity.ExportJobRequestedEvent{
			JobID:    item.ID,
			TenantID: item.TenantID,
		})
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}
