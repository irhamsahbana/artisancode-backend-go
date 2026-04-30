package core

import (
	"testing"

	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

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

	require.Equal(t, shiftName, row[3])
	require.Equal(t, "15", row[7])
	require.Equal(t, selfieURL, row[16])
	require.Equal(t, "https://www.google.com/maps?q=-0.8762091,119.8875055", row[17])
}
