package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttendanceExportHeaders_LocalizedToIndonesian(t *testing.T) {
	headers := attendanceExportHeaders("id")

	require.Equal(t, "No Karyawan", headers[0])
	require.Equal(t, "Link Lokasi", headers[len(headers)-1])
}
