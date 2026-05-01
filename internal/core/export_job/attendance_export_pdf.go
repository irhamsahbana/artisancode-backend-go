package core

import (
	"bytes"
	"fmt"
	"strings"

	"codebase-app/internal/entity/coreentity"

	"github.com/jung-kurt/gofpdf"
)

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

func truncateForPDF(value string, width float64) string {
	limit := int(width * 0.9)
	if limit < 4 || len(value) <= limit {
		return value
	}

	return value[:limit-3] + "..."
}
