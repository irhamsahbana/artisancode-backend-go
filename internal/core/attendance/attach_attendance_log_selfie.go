package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) attachAttendanceLogSelfie(ctx context.Context, item *coreentity.AttendanceLog) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:attendance:attach_attendance_log_selfie:attachAttendanceLogSelfie",
	)
	defer span.End()

	if item == nil {
		return nil
	}

	link, err := c.storageRepo.GetFileLink(ctx, coreentity.StorageFileLinkFilter{
		TenantID:     item.TenantID,
		ResourceType: "attendance_log",
		ResourceID:   item.ID,
		FieldName:    "selfie",
	})
	if err != nil {
		return err
	}
	if link == nil {
		return nil
	}

	file, err := c.storageRepo.GetFile(ctx, coreentity.FileFilter{
		TenantID: item.TenantID,
		ID:       link.StorageFileID,
	})
	if err != nil {
		return err
	}

	item.SelfieFileID = &file.ID

	url, err := c.s3.GetFileURL(ctx, coreentity.FileFilter{
		TenantID: file.TenantID,
		Filename: file.Filename,
	})
	if err != nil {
		return err
	}
	item.SelfieURL = &url

	return nil
}
