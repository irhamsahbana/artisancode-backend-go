package core

import (
	"bytes"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

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
	require.NoError(t, err)

	file, err := excelize.OpenReader(bytes.NewReader(content))
	require.NoError(t, err)
	defer file.Close()

	require.Equal(t, []string{"Attendance Logs", "Employee Summary"}, file.GetSheetList())

	got, err := file.GetCellValue("Employee Summary", "A2")
	require.NoError(t, err)
	require.Equal(t, "EMP-001", got)
}
