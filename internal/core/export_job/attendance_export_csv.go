package core

import (
	"bytes"
	"encoding/csv"

	"codebase-app/internal/entity/coreentity"
)

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
