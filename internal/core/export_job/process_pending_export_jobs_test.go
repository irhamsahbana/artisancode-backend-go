package core

import (
	"bytes"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"github.com/xuri/excelize/v2"
)

func TestAttendanceExportHeaders_LocalizedToIndonesian(t *testing.T) {
	headers := attendanceExportHeaders("id")

	if headers[0] != "No Karyawan" {
		t.Fatalf("expected Indonesian employee number header, got %q", headers[0])
	}

	if headers[len(headers)-1] != "Link Lokasi" {
		t.Fatalf("expected Indonesian location link header, got %q", headers[len(headers)-1])
	}
}

func TestAttendanceExportRow_UsesSignedPhotoAndGoogleMapsLink(t *testing.T) {
	latitude := -0.8762091
	longitude := 119.8875055
	selfieURL := "https://signed.example.com/selfie"
	shiftName := "Morning Shift"
	shiftTimezone := "Asia/Makassar"
	shiftStartTime := "08:00"
	shiftEndTime := "17:00"
	shiftGraceMins := 15

	row := attendanceExportRow(coreentity.AttendanceLog{
		EmployeeNo:     "EMP-001",
		EmployeeName:   "Budi",
		ShiftName:      &shiftName,
		ShiftTimezone:  &shiftTimezone,
		ShiftStartTime: &shiftStartTime,
		ShiftEndTime:   &shiftEndTime,
		ShiftGraceMins: &shiftGraceMins,
		SelfieURL:      &selfieURL,
		Latitude:       &latitude,
		Longitude:      &longitude,
	})

	if got := row[3]; got != shiftName {
		t.Fatalf("expected shift name in export row, got %q", got)
	}

	if got := row[7]; got != "15" {
		t.Fatalf("expected shift grace minutes in export row, got %q", got)
	}

	if got := row[16]; got != selfieURL {
		t.Fatalf("expected signed selfie url in photo proof column, got %q", got)
	}

	expectedMapsURL := "https://www.google.com/maps?q=-0.8762091,119.8875055"
	if got := row[17]; got != expectedMapsURL {
		t.Fatalf("expected Google Maps url %q, got %q", expectedMapsURL, got)
	}
}

func TestBuildAttendanceEmployeeSummaries_AggregatesPerEmployee(t *testing.T) {
	shiftStart := "08:00"
	shiftTimezone := "Asia/Makassar"
	shiftGraceMins := 15
	selfieURL := "https://signed.example.com/selfie"

	summaries := buildAttendanceEmployeeSummaries([]coreentity.AttendanceLog{
		{
			EmployeeID:     "emp-1",
			EmployeeNo:     "EMP-001",
			EmployeeName:   "Budi",
			AttendanceDate: "2026-04-18",
			Type:           common.AttendanceTypeCheckIn,
			LoggedAt:       "2026-04-18T08:30:00+08:00",
			ShiftStartTime: &shiftStart,
			ShiftTimezone:  &shiftTimezone,
			ShiftGraceMins: &shiftGraceMins,
			SelfieURL:      &selfieURL,
		},
		{
			EmployeeID:     "emp-1",
			EmployeeNo:     "EMP-001",
			EmployeeName:   "Budi",
			AttendanceDate: "2026-04-18",
			Type:           common.AttendanceTypeCheckOut,
			LoggedAt:       "2026-04-18T17:01:00+08:00",
		},
		{
			EmployeeID:     "emp-1",
			EmployeeNo:     "EMP-001",
			EmployeeName:   "Budi",
			AttendanceDate: "2026-04-19",
			Type:           common.AttendanceTypeCheckIn,
			LoggedAt:       "2026-04-19T08:05:00+08:00",
			ShiftStartTime: &shiftStart,
			ShiftTimezone:  &shiftTimezone,
			ShiftGraceMins: &shiftGraceMins,
		},
		{
			EmployeeID:     "emp-2",
			EmployeeNo:     "EMP-002",
			EmployeeName:   "Siti",
			AttendanceDate: "2026-04-18",
			Type:           common.AttendanceTypeCheckOut,
			LoggedAt:       "2026-04-18T17:05:00+08:00",
		},
	})

	if len(summaries) != 2 {
		t.Fatalf("expected 2 employee summaries, got %d", len(summaries))
	}

	first := summaries[0]
	if first.EmployeeNo != "EMP-001" {
		t.Fatalf("expected first summary for EMP-001, got %q", first.EmployeeNo)
	}
	if first.AttendanceDays != 2 {
		t.Fatalf("expected 2 attendance days, got %d", first.AttendanceDays)
	}
	if first.TotalLogs != 3 {
		t.Fatalf("expected 3 total logs, got %d", first.TotalLogs)
	}
	if first.LateCheckInCount != 1 {
		t.Fatalf("expected 1 late check-in day, got %d", first.LateCheckInCount)
	}
	if first.MissingCheckOutDays != 1 {
		t.Fatalf("expected 1 missing check-out day, got %d", first.MissingCheckOutDays)
	}
	if first.WithPhotoCount != 1 || first.WithoutPhotoCount != 2 {
		t.Fatalf("expected photo counts 1/2, got %d/%d", first.WithPhotoCount, first.WithoutPhotoCount)
	}

	second := summaries[1]
	if second.MissingCheckInDays != 1 {
		t.Fatalf("expected 1 missing check-in day for second employee, got %d", second.MissingCheckInDays)
	}
}

func TestBuildAttendanceXLSX_AddsEmployeeSummarySheet(t *testing.T) {
	content, err := buildAttendanceXLSX([]coreentity.AttendanceLog{
		{
			EmployeeID:     "emp-1",
			EmployeeNo:     "EMP-001",
			EmployeeName:   "Budi",
			AttendanceDate: "2026-04-18",
			Type:           common.AttendanceTypeCheckIn,
			Status:         common.AttendanceStatusRecorded,
			Source:         common.AttendanceSourceMobile,
			LoggedAt:       "2026-04-18T08:00:00+08:00",
		},
	}, "en")
	if err != nil {
		t.Fatalf("expected xlsx generation to succeed, got error %v", err)
	}

	file, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("expected valid xlsx content, got error %v", err)
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) != 2 {
		t.Fatalf("expected 2 sheets, got %d", len(sheets))
	}
	if sheets[0] != "Attendance Logs" {
		t.Fatalf("expected first sheet Attendance Logs, got %q", sheets[0])
	}
	if sheets[1] != "Employee Summary" {
		t.Fatalf("expected second sheet Employee Summary, got %q", sheets[1])
	}

	if got, err := file.GetCellValue("Employee Summary", "A2"); err != nil || got != "EMP-001" {
		t.Fatalf("expected employee summary row to contain employee number, got value=%q err=%v", got, err)
	}
}
