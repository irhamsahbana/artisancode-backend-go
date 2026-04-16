package common

type ModerationStatus string

const (
	ModerationStatusCreated  ModerationStatus = "created"
	ModerationStatusApproved ModerationStatus = "approved"
	ModerationStatusRejected ModerationStatus = "rejected"
)

type S3Folder string

const (
	S3FolderAttendanceReports S3Folder = "attandance-reports"
	S3FolderAttendanceFace    S3Folder = "attendance-face"
)

type AttendanceType string

const (
	AttendanceTypeCheckIn  AttendanceType = "check_in"
	AttendanceTypeCheckOut AttendanceType = "check_out"
)

type AttendanceSource string

const (
	AttendanceSourceWeb    AttendanceSource = "web"
	AttendanceSourceMobile AttendanceSource = "mobile"
)

type AttendanceStatus string

const (
	AttendanceStatusRecorded AttendanceStatus = "recorded"
)

type AttendanceTodayStatus string

const (
	AttendanceTodayStatusNotCheckedIn AttendanceTodayStatus = "not_checked_in"
	AttendanceTodayStatusCheckedIn    AttendanceTodayStatus = "checked_in"
	AttendanceTodayStatusCheckedOut   AttendanceTodayStatus = "checked_out"
)

type FileStatus string

const (
	FileStatusPending  FileStatus = "pending"
	FileStatusAttached FileStatus = "attached"
	FileStatusDeleted  FileStatus = "deleted"
	FileStatusFailed   FileStatus = "failed"
)

const (
	LogKeyPayload string = "payload"
)
