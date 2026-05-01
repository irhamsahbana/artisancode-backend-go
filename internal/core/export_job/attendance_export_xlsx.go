package core

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"

	"github.com/xuri/excelize/v2"
)

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

	file.SetPanes(
		logSheetName,
		&excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2"},
	)

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
