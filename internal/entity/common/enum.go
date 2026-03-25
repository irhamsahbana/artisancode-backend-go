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
)

const (
	LogKeyPayload string = "payload"
)
