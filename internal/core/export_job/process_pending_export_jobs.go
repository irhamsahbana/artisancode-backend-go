package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

type exportFilters struct {
	Language       string  `json:"language"`
	Q              string  `json:"q"`
	EmployeeID     *string `json:"employee_id"`
	Type           string  `json:"type"`
	Source         string  `json:"source"`
	Status         string  `json:"status"`
	SelfieStatus   string  `json:"selfie_status"`
	OrgUnitID      *string `json:"org_unit_id"`
	BranchID       *string `json:"branch_id"`
	WorkLocationID *string `json:"work_location_id"`
	ExceptionType  string  `json:"exception_type"`
	AttendanceDay  *string `json:"attendance_date"`
	DateFrom       *string `json:"date_from"`
	DateTo         *string `json:"date_to"`
}

func (c *exportJobCore) ProcessPendingExportJobs(
	ctx context.Context,
	limit int,
) (*coreentity.ExportJobProcessResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_pending_export_jobs:ProcessPendingExportJobs")
	defer span.End()

	if limit < 1 {
		limit = 1
	}

	result := &coreentity.ExportJobProcessResult{}
	for index := 0; index < limit; index++ {
		item, err := c.repo.ClaimPendingExportJob(ctx)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return result, nil
		}

		err = c.processClaimedExport(ctx, *item)
		if err != nil {
			errorMessage := err.Error()
			_ = c.repo.UpdateExportJob(ctx, coreentity.ExportJobUpdate{
				TenantID:     item.TenantID,
				ID:           item.ID,
				Status:       coreentity.ExportJobStatusFailed,
				ErrorMessage: &errorMessage,
				StartedAt:    item.StartedAt,
			})
			return nil, err
		}

		result.Processed++
	}

	return result, nil
}

func (c *exportJobCore) processClaimedExport(ctx context.Context, item coreentity.ExportJob) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_pending_export_jobs:processClaimedExport")
	defer span.End()

	return c.processJob(ctx, item)
}

func (c *exportJobCore) storeGeneratedFile(
	ctx context.Context,
	item coreentity.ExportJob,
	originalFilename, contentType string,
	content []byte,
) (string, string, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_pending_export_jobs:storeGeneratedFile")
	defer span.End()

	filename := buildAttendanceReportObjectKey(item.TenantID, item.RequestedBy, item.ID, originalFilename)
	sizeBytes := int64(len(content))
	file, err := c.storageRepo.CreateFile(ctx, coreentity.File{
		TenantID:         item.TenantID,
		CreatedBy:        item.RequestedBy,
		Status:           common.FileStatusPending,
		Folder:           common.S3FolderAttendanceReports,
		Filename:         filename,
		OriginalFilename: &originalFilename,
		ContentType:      &contentType,
		SizeBytes:        &sizeBytes,
		IsPublic:         false,
	})
	if err != nil {
		return "", "", err
	}

	_, err = c.s3.UploadBytes(ctx, &coreentity.UploadBytesReq{
		Filename:    filename,
		ContentType: contentType,
		Body:        content,
		IsPublic:    false,
	})
	if err != nil {
		return "", "", err
	}

	_, err = c.storageRepo.CreateFileLink(ctx, coreentity.CreateStorageFileLinkReq{
		TenantID:      item.TenantID,
		StorageFileID: file.ID,
		ResourceType:  "export_job",
		ResourceID:    item.ID,
		FieldName:     "report_file",
		SortOrder:     1,
	})
	if err != nil {
		return "", "", err
	}

	err = c.storageRepo.MarkFileAttached(ctx, item.TenantID, file.ID)
	if err != nil {
		return "", "", err
	}

	if file.ExpiresAt == nil {
		return file.ID, time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339), nil
	}

	return file.ID, *file.ExpiresAt, nil
}

func (c *exportJobCore) attachDownloadURL(ctx context.Context, item *coreentity.ExportJob, now time.Time) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:export_job:process_pending_export_jobs:attachDownloadURL")
	defer span.End()

	if item.Status == coreentity.ExportJobStatusCompleted && item.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.RFC3339, *item.ExpiresAt)
		if err == nil && now.After(expiresAt) {
			item.Status = coreentity.ExportJobStatusExpired
			item.DownloadURL = nil
			return nil
		}
	}

	if item.Status != coreentity.ExportJobStatusCompleted || item.FileID == nil {
		return nil
	}

	file, err := c.storageRepo.GetFile(ctx, coreentity.FileFilter{
		TenantID: item.TenantID,
		ID:       *item.FileID,
	})
	if err != nil {
		return err
	}

	url, err := c.s3.GetFileURL(ctx, coreentity.FileFilter{
		TenantID: item.TenantID,
		Filename: file.Filename,
	})
	if err != nil {
		return err
	}

	item.DownloadURL = &url
	item.FileName = file.OriginalFilename
	return nil
}

func buildAttendanceReportObjectKey(tenantID, userID, exportID, filename string) string {
	return fmt.Sprintf(
		"private/tenants/%s/%s/users/%s/%s-%d-%s",
		tenantID,
		strings.TrimPrefix(string(common.S3FolderAttendanceReports), "/"),
		userID,
		exportID,
		time.Now().UTC().UnixNano(),
		strings.ReplaceAll(filename, " ", "-"),
	)
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func attendanceTypePtr(value string) *common.AttendanceType {
	if value == "" {
		return nil
	}
	result := common.AttendanceType(value)
	return &result
}

func attendanceSourcePtr(value string) *common.AttendanceSource {
	if value == "" {
		return nil
	}
	result := common.AttendanceSource(value)
	return &result
}

func attendanceStatusPtr(value string) *common.AttendanceStatus {
	if value == "" {
		return nil
	}
	result := common.AttendanceStatus(value)
	return &result
}
