package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *exportJobCore) ProcessExportJob(ctx context.Context, filter coreentity.ExportJobDetailFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_export_job:ProcessExportJob")
	defer span.End()

	item, err := c.repo.GetExportJob(ctx, filter)
	if err != nil {
		return err
	}

	if item.Status != coreentity.ExportJobStatusPending {
		return nil
	}

	startedAt := time.Now().UTC().Format(time.RFC3339)
	err = c.repo.UpdateExportJob(ctx, coreentity.ExportJobUpdate{
		TenantID:  item.TenantID,
		ID:        item.ID,
		Status:    coreentity.ExportJobStatusProcessing,
		StartedAt: &startedAt,
	})
	if err != nil {
		return err
	}
	item.Status = coreentity.ExportJobStatusProcessing
	item.StartedAt = &startedAt

	return c.processJob(ctx, *item)
}

func (c *exportJobCore) processJob(ctx context.Context, item coreentity.ExportJob) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_export_job:processJob")
	defer span.End()

	if item.ProcessorKey != "attendance_logs" {
		return fmt.Errorf("unsupported export job processor: %s", item.ProcessorKey)
	}

	var filters exportFilters
	err := json.Unmarshal([]byte(item.ParamsJSON), &filters)
	if err != nil {
		return err
	}

	logs, err := c.attendanceRepo.GetAttendanceLogsAll(ctx, coreentity.AttendanceLogListFilter{
		UserCtx:        common.UserContext{UserID: item.RequestedBy, TenantID: item.TenantID, Roles: []string{"owner"}},
		TenantID:       item.TenantID,
		EmployeeID:     filters.EmployeeID,
		Q:              filters.Q,
		Type:           attendanceTypePtr(filters.Type),
		Source:         attendanceSourcePtr(filters.Source),
		Status:         attendanceStatusPtr(filters.Status),
		SelfieStatus:   stringPtr(filters.SelfieStatus),
		OrgUnitID:      filters.OrgUnitID,
		BranchID:       filters.BranchID,
		WorkLocationID: filters.WorkLocationID,
		ExceptionType:  stringPtr(filters.ExceptionType),
		AttendanceDay:  filters.AttendanceDay,
		DateFrom:       filters.DateFrom,
		DateTo:         filters.DateTo,
		Page:           1,
		Paginate:       1000000,
	})
	if err != nil {
		return err
	}

	content, filename, contentType, err := generateAttendanceReportFile(item.Format, logs, filters)
	if err != nil {
		return err
	}

	fileID, expiresAt, err := c.storeGeneratedFile(ctx, item, filename, contentType, content)
	if err != nil {
		return err
	}

	completedAt := time.Now().UTC().Format(time.RFC3339)
	return c.repo.UpdateExportJob(ctx, coreentity.ExportJobUpdate{
		TenantID:    item.TenantID,
		ID:          item.ID,
		Status:      coreentity.ExportJobStatusCompleted,
		FileID:      &fileID,
		StartedAt:   item.StartedAt,
		CompletedAt: &completedAt,
		ExpiresAt:   &expiresAt,
	})
}
