package core

import (
	"fmt"
	"time"

	"codebase-app/internal/entity/coreentity"
)

func generateAttendanceReportFile(
	format coreentity.ExportJobFormat,
	logs []coreentity.AttendanceLog,
	filters exportFilters,
) ([]byte, string, string, error) {
	filenameBase := fmt.Sprintf("attendance-report-%s", time.Now().UTC().Format("20060102-150405"))

	switch format {
	// case coreentity.ExportJobFormatCSV:
	// 	content, err := buildAttendanceCSV(logs, filters.Language)
	// 	return content, filenameBase + ".csv", "text/csv", err
	case coreentity.ExportJobFormatXLSX:
		content, err := buildAttendanceXLSX(logs, filters.Language)
		return content, filenameBase + ".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", err
	// case coreentity.ExportJobFormatPDF:
	// 	content, err := buildAttendancePDF(logs, filters)
	// 	return content, filenameBase + ".pdf", "application/pdf", err
	default:
		return nil, "", "", fmt.Errorf("unsupported export format: %s", format)
	}
}
