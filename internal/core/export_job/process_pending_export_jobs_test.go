package core

import (
	"testing"

	"codebase-app/internal/entity/coreentity"
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
