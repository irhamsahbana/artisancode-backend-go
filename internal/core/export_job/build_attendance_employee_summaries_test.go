package core

import (
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
)

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

	require.Len(t, summaries, 2)

	first := summaries[0]
	require.Equal(t, "EMP-001", first.EmployeeNo)
	require.Equal(t, 2, first.AttendanceDays)
	require.Equal(t, 3, first.TotalLogs)
	require.Equal(t, 1, first.LateCheckInCount)
	require.Equal(t, 1, first.MissingCheckOutDays)
	require.Equal(t, 1, first.WithPhotoCount)
	require.Equal(t, 2, first.WithoutPhotoCount)

	require.Equal(t, 1, summaries[1].MissingCheckInDays)
}
