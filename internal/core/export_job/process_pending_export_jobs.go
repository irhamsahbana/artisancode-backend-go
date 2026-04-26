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
	"sort"
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

func generateAttendanceReportFile(
	format coreentity.ExportJobFormat,
	logs []coreentity.AttendanceLog,
	filters exportFilters,
) ([]byte, string, string, error) {
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

	logSheetName := localizeExportText(language, "Attendance Logs", "Log Kehadiran")
	file.SetSheetName("Sheet1", logSheetName)

	headers := attendanceExportHeaders(language)
	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, 1)
		file.SetCellValue(logSheetName, cell, header)
	}

	for rowIndex, item := range logs {
		for columnIndex, value := range attendanceExportRow(item) {
			cell, _ := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+2)
			file.SetCellValue(logSheetName, cell, value)
		}
	}

	file.SetPanes(logSheetName, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2"})

	summarySheetName := localizeExportText(language, "Employee Summary", "Ringkasan Karyawan")
	file.NewSheet(summarySheetName)

	summaryHeaders := attendanceEmployeeSummaryHeaders(language)
	for index, header := range summaryHeaders {
		cell, _ := excelize.CoordinatesToCellName(index+1, 1)
		file.SetCellValue(summarySheetName, cell, header)
	}

	for rowIndex, item := range buildAttendanceEmployeeSummaries(logs) {
		for columnIndex, value := range attendanceEmployeeSummaryRow(item) {
			cell, _ := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+2)
			file.SetCellValue(summarySheetName, cell, value)
		}
	}

	file.SetPanes(
		summarySheetName,
		&excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2"},
	)
	logSheetIndex, err := file.GetSheetIndex(logSheetName)
	if err == nil {
		file.SetActiveSheet(logSheetIndex)
	}
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
	pdf.CellFormat(
		0,
		8,
		localizeExportText(filters.Language, "Attendance Report", "Laporan Kehadiran"),
		"",
		1,
		"L",
		false,
		0,
		"",
	)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5, buildAttendanceFilterSummary(filters), "", "L", false)
	pdf.Ln(2)

	headers := attendanceExportHeaders(filters.Language)
	widths := attendancePDFColumnWidths(len(headers))

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
		localizeExportText(language, "Shift Name", "Nama Shift"),
		localizeExportText(language, "Shift Timezone", "Zona Waktu Shift"),
		localizeExportText(language, "Shift Start", "Mulai Shift"),
		localizeExportText(language, "Shift End", "Selesai Shift"),
		localizeExportText(language, "Shift Grace (Minutes)", "Toleransi Shift (Menit)"),
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
		stringValue(item.ShiftName),
		stringValue(item.ShiftTimezone),
		stringValue(item.ShiftStartTime),
		stringValue(item.ShiftEndTime),
		intValue(item.ShiftGraceMins),
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

type attendanceEmployeeSummary struct {
	EmployeeID          string
	EmployeeNo          string
	EmployeeName        string
	AttendanceDays      int
	TotalLogs           int
	CheckInCount        int
	CheckOutCount       int
	LateCheckInCount    int
	MissingCheckInDays  int
	MissingCheckOutDays int
	WithPhotoCount      int
	WithoutPhotoCount   int
	FirstAttendanceDate string
	LastAttendanceDate  string
	LastLoggedAt        string
}

type attendanceEmployeeDaySummary struct {
	HasCheckIn       bool
	HasCheckOut      bool
	HasLateCheckIn   bool
	EarliestCheckIn  *coreentity.AttendanceLog
	LatestAttendance *coreentity.AttendanceLog
}

func attendanceEmployeeSummaryHeaders(language string) []string {
	return []string{
		localizeExportText(language, "Employee No", "No Karyawan"),
		localizeExportText(language, "Employee Name", "Nama Karyawan"),
		localizeExportText(language, "Attendance Days", "Hari Kehadiran"),
		localizeExportText(language, "Total Logs", "Total Log"),
		localizeExportText(language, "Check In Count", "Jumlah Check-in"),
		localizeExportText(language, "Check Out Count", "Jumlah Check-out"),
		localizeExportText(language, "Late Check In Days", "Hari Terlambat Check-in"),
		localizeExportText(language, "Missing Check In Days", "Hari Tanpa Check-in"),
		localizeExportText(language, "Missing Check Out Days", "Hari Tanpa Check-out"),
		localizeExportText(language, "With Photo", "Dengan Foto"),
		localizeExportText(language, "Without Photo", "Tanpa Foto"),
		localizeExportText(language, "First Attendance Date", "Tanggal Kehadiran Pertama"),
		localizeExportText(language, "Last Attendance Date", "Tanggal Kehadiran Terakhir"),
		localizeExportText(language, "Last Logged At", "Waktu Log Terakhir"),
	}
}

func attendanceEmployeeSummaryRow(item attendanceEmployeeSummary) []string {
	return []string{
		item.EmployeeNo,
		item.EmployeeName,
		fmt.Sprintf("%d", item.AttendanceDays),
		fmt.Sprintf("%d", item.TotalLogs),
		fmt.Sprintf("%d", item.CheckInCount),
		fmt.Sprintf("%d", item.CheckOutCount),
		fmt.Sprintf("%d", item.LateCheckInCount),
		fmt.Sprintf("%d", item.MissingCheckInDays),
		fmt.Sprintf("%d", item.MissingCheckOutDays),
		fmt.Sprintf("%d", item.WithPhotoCount),
		fmt.Sprintf("%d", item.WithoutPhotoCount),
		item.FirstAttendanceDate,
		item.LastAttendanceDate,
		formatAttendanceTimestamp(item.LastLoggedAt),
	}
}

func buildAttendanceEmployeeSummaries(logs []coreentity.AttendanceLog) []attendanceEmployeeSummary {
	type employeeAggregate struct {
		summary attendanceEmployeeSummary
		days    map[string]*attendanceEmployeeDaySummary
	}

	aggregates := make(map[string]*employeeAggregate)

	for _, item := range logs {
		key := item.EmployeeID
		if key == "" {
			key = item.EmployeeNo + ":" + item.EmployeeName
		}

		aggregate, exists := aggregates[key]
		if !exists {
			aggregate = &employeeAggregate{
				summary: attendanceEmployeeSummary{
					EmployeeID:   item.EmployeeID,
					EmployeeNo:   item.EmployeeNo,
					EmployeeName: item.EmployeeName,
				},
				days: make(map[string]*attendanceEmployeeDaySummary),
			}
			aggregates[key] = aggregate
		}

		aggregate.summary.TotalLogs++
		if item.Type == common.AttendanceTypeCheckIn {
			aggregate.summary.CheckInCount++
		}
		if item.Type == common.AttendanceTypeCheckOut {
			aggregate.summary.CheckOutCount++
		}
		if item.SelfieURL != nil && *item.SelfieURL != "" {
			aggregate.summary.WithPhotoCount++
		} else {
			aggregate.summary.WithoutPhotoCount++
		}
		if aggregate.summary.FirstAttendanceDate == "" || item.AttendanceDate < aggregate.summary.FirstAttendanceDate {
			aggregate.summary.FirstAttendanceDate = item.AttendanceDate
		}
		if item.AttendanceDate > aggregate.summary.LastAttendanceDate {
			aggregate.summary.LastAttendanceDate = item.AttendanceDate
		}

		loggedAt, err := time.Parse(time.RFC3339, item.LoggedAt)
		if err == nil {
			lastLoggedAt, lastErr := time.Parse(time.RFC3339, aggregate.summary.LastLoggedAt)
			if aggregate.summary.LastLoggedAt == "" || lastErr != nil || loggedAt.After(lastLoggedAt) {
				aggregate.summary.LastLoggedAt = item.LoggedAt
			}
		}

		daySummary, exists := aggregate.days[item.AttendanceDate]
		if !exists {
			daySummary = &attendanceEmployeeDaySummary{}
			aggregate.days[item.AttendanceDate] = daySummary
		}

		updateAttendanceEmployeeDaySummary(daySummary, item)
	}

	summaries := make([]attendanceEmployeeSummary, 0, len(aggregates))
	for _, aggregate := range aggregates {
		aggregate.summary.AttendanceDays = len(aggregate.days)
		for _, daySummary := range aggregate.days {
			if daySummary.HasLateCheckIn {
				aggregate.summary.LateCheckInCount++
			}
			if daySummary.HasCheckOut && !daySummary.HasCheckIn {
				aggregate.summary.MissingCheckInDays++
			}
			if daySummary.HasCheckIn && !daySummary.HasCheckOut {
				aggregate.summary.MissingCheckOutDays++
			}
		}

		summaries = append(summaries, aggregate.summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].EmployeeNo != summaries[j].EmployeeNo {
			return summaries[i].EmployeeNo < summaries[j].EmployeeNo
		}
		return summaries[i].EmployeeName < summaries[j].EmployeeName
	})

	return summaries
}

func updateAttendanceEmployeeDaySummary(daySummary *attendanceEmployeeDaySummary, item coreentity.AttendanceLog) {
	if item.Type == common.AttendanceTypeCheckIn {
		daySummary.HasCheckIn = true
		if daySummary.EarliestCheckIn == nil ||
			compareAttendanceLoggedAt(item.LoggedAt, daySummary.EarliestCheckIn.LoggedAt) < 0 {
			checkIn := item
			daySummary.EarliestCheckIn = &checkIn
			daySummary.HasLateCheckIn = isLateCheckIn(item)
		}
	}

	if item.Type == common.AttendanceTypeCheckOut {
		daySummary.HasCheckOut = true
	}
}

func compareAttendanceLoggedAt(left, right string) int {
	leftTime, leftErr := time.Parse(time.RFC3339, left)
	rightTime, rightErr := time.Parse(time.RFC3339, right)

	if leftErr != nil || rightErr != nil {
		return strings.Compare(left, right)
	}

	if leftTime.Before(rightTime) {
		return -1
	}
	if leftTime.After(rightTime) {
		return 1
	}
	return 0
}

func isLateCheckIn(item coreentity.AttendanceLog) bool {
	if item.Type != common.AttendanceTypeCheckIn || item.ShiftStartTime == nil || *item.ShiftStartTime == "" {
		return false
	}

	loggedAt, err := time.Parse(time.RFC3339, item.LoggedAt)
	if err != nil {
		return false
	}

	location := time.UTC
	if item.ShiftTimezone != nil && *item.ShiftTimezone != "" {
		timezoneLocation, locationErr := time.LoadLocation(*item.ShiftTimezone)
		if locationErr == nil {
			location = timezoneLocation
		}
	}

	shiftStart, err := parseAttendanceClock(*item.ShiftStartTime)
	if err != nil {
		return false
	}

	graceMinutes := 0
	if item.ShiftGraceMins != nil {
		graceMinutes = *item.ShiftGraceMins
	}

	localLoggedAt := loggedAt.In(location)
	shiftStartAt := time.Date(
		localLoggedAt.Year(),
		localLoggedAt.Month(),
		localLoggedAt.Day(),
		shiftStart.Hour(),
		shiftStart.Minute(),
		shiftStart.Second(),
		0,
		location,
	).Add(time.Duration(graceMinutes) * time.Minute)

	return localLoggedAt.After(shiftStartAt)
}

func parseAttendanceClock(value string) (time.Time, error) {
	layouts := []string{"15:04:05", "15:04"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid attendance clock: %s", value)
}

func intValue(value *int) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%d", *value)
}

func attendancePDFColumnWidths(columnCount int) []float64 {
	widths := []float64{20, 32, 24, 24, 24, 18, 18, 18, 16, 16, 16, 28, 38, 22, 24, 30, 48, 44}
	if columnCount <= len(widths) {
		return widths[:columnCount]
	}

	for len(widths) < columnCount {
		widths = append(widths, 24)
	}

	return widths
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
