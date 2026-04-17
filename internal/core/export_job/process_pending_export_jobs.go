package core

import (
	"bytes"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
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

func (c *exportJobCore) ProcessPendingExportJobs(ctx context.Context, limit int) (*coreentity.ExportJobProcessResult, error) {
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

func generateAttendanceReportFile(format coreentity.ExportJobFormat, logs []coreentity.AttendanceLog, filters exportFilters) ([]byte, string, string, error) {
	filenameBase := fmt.Sprintf("attendance-report-%s", time.Now().UTC().Format("20060102-150405"))

	switch format {
	case coreentity.ExportJobFormatCSV:
		content, err := buildAttendanceCSV(logs, filters.Language)
		return content, filenameBase + ".csv", "text/csv", err
	case coreentity.ExportJobFormatXLSX:
		content, err := buildAttendanceXLSX(logs, filters.Language)
		return content, filenameBase + ".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", err
	case coreentity.ExportJobFormatPDF:
		content, err := buildAttendancePDF(logs, filters)
		return content, filenameBase + ".pdf", "application/pdf", err
	default:
		return nil, "", "", fmt.Errorf("unsupported export format: %s", format)
	}
}

func buildAttendanceCSV(logs []coreentity.AttendanceLog, language string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	writer := csv.NewWriter(buffer)

	err := writer.Write(attendanceExportHeaders(language))
	if err != nil {
		return nil, err
	}

	for _, item := range logs {
		err = writer.Write(attendanceExportRow(item))
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func buildAttendanceXLSX(logs []coreentity.AttendanceLog, language string) ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()

	sheetName := "Attendance Logs"
	file.SetSheetName("Sheet1", sheetName)

	headers := attendanceExportHeaders(language)
	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for rowIndex, item := range logs {
		for columnIndex, value := range attendanceExportRow(item) {
			cell, _ := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+2)
			file.SetCellValue(sheetName, cell, value)
		}
	}

	file.SetPanes(sheetName, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2"})
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func buildAttendancePDF(logs []coreentity.AttendanceLog, filters exportFilters) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(8, 8, 8)
	pdf.SetAutoPageBreak(true, 8)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, localizeExportText(filters.Language, "Attendance Report", "Laporan Kehadiran"), "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5, buildAttendanceFilterSummary(filters), "", "L", false)
	pdf.Ln(2)

	headers := attendanceExportHeaders(filters.Language)
	widths := []float64{20, 32, 24, 16, 16, 16, 28, 38, 22, 24, 30, 48, 44}

	pdf.SetFont("Arial", "B", 7)
	for index, header := range headers {
		pdf.CellFormat(widths[index], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 6)
	for _, item := range logs {
		row := attendanceExportRow(item)
		for index, value := range row {
			pdf.CellFormat(widths[index], 6, truncateForPDF(value, widths[index]), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	buffer := bytes.NewBuffer(nil)
	err := pdf.Output(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func attendanceExportHeaders(language string) []string {
	return []string{
		localizeExportText(language, "Employee No", "No Karyawan"),
		localizeExportText(language, "Employee Name", "Nama Karyawan"),
		localizeExportText(language, "Attendance Date", "Tanggal Kehadiran"),
		localizeExportText(language, "Type", "Tipe"),
		localizeExportText(language, "Source", "Sumber"),
		localizeExportText(language, "Status", "Status"),
		localizeExportText(language, "Logged At", "Waktu Check"),
		localizeExportText(language, "Address", "Alamat"),
		localizeExportText(language, "Device ID", "ID Perangkat"),
		localizeExportText(language, "Device Name", "Nama Perangkat"),
		localizeExportText(language, "Notes", "Catatan"),
		localizeExportText(language, "Photo Proof", "Bukti Foto"),
		localizeExportText(language, "Location Link", "Link Lokasi"),
	}
}

func attendanceExportRow(item coreentity.AttendanceLog) []string {
	return []string{
		item.EmployeeNo,
		item.EmployeeName,
		item.AttendanceDate,
		string(item.Type),
		string(item.Source),
		string(item.Status),
		formatAttendanceTimestamp(item.LoggedAt),
		stringValue(item.Address),
		stringValue(item.DeviceID),
		stringValue(item.DeviceName),
		stringValue(item.Notes),
		stringValue(item.SelfieURL),
		buildGoogleMapsURL(item.Latitude, item.Longitude),
	}
}

func buildAttendanceFilterSummary(filters exportFilters) string {
	parts := []string{
		localizeExportText(filters.Language, "Filters:", "Filter:"),
	}
	if filters.DateFrom != nil || filters.DateTo != nil {
		parts = append(parts, fmt.Sprintf("date=%s to %s", stringValue(filters.DateFrom), stringValue(filters.DateTo)))
	}
	if filters.AttendanceDay != nil {
		parts = append(parts, "attendance_date="+*filters.AttendanceDay)
	}
	if filters.Type != "" {
		parts = append(parts, "type="+filters.Type)
	}
	if filters.Source != "" {
		parts = append(parts, "source="+filters.Source)
	}
	if filters.Status != "" {
		parts = append(parts, "status="+filters.Status)
	}
	if filters.ExceptionType != "" {
		parts = append(parts, "exception="+filters.ExceptionType)
	}
	if filters.SelfieStatus != "" {
		parts = append(parts, "photo="+filters.SelfieStatus)
	}
	if filters.Q != "" {
		parts = append(parts, "search="+filters.Q)
	}

	return strings.Join(parts, " | ")
}

func localizeExportText(language, englishText, indonesianText string) string {
	if errmsg.ResolveLanguage(language) == errmsg.LanguageEnglish {
		return englishText
	}

	return indonesianText
}

func buildGoogleMapsURL(latitude, longitude *float64) string {
	if latitude == nil || longitude == nil {
		return ""
	}

	return fmt.Sprintf("https://www.google.com/maps?q=%.7f,%.7f", *latitude, *longitude)
}

func (c *exportJobCore) storeGeneratedFile(ctx context.Context, item coreentity.ExportJob, originalFilename, contentType string, content []byte) (string, string, error) {
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

func formatAttendanceTimestamp(value string) string {
	if value == "" {
		return ""
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}

	return parsed.Format("2006-01-02 15:04:05 MST")
}

func truncateForPDF(value string, width float64) string {
	limit := int(width * 0.9)
	if limit < 4 || len(value) <= limit {
		return value
	}

	return value[:limit-3] + "..."
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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
