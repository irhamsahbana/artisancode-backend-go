package coreentity

import "codebase-app/internal/entity/common"

type ExportJobFormat string

const (
	ExportJobFormatCSV  ExportJobFormat = "csv"
	ExportJobFormatXLSX ExportJobFormat = "xlsx"
	ExportJobFormatPDF  ExportJobFormat = "pdf"
)

type ExportJobStatus string

const (
	ExportJobStatusPending    ExportJobStatus = "pending"
	ExportJobStatusProcessing ExportJobStatus = "processing"
	ExportJobStatusCompleted  ExportJobStatus = "completed"
	ExportJobStatusFailed     ExportJobStatus = "failed"
	ExportJobStatusExpired    ExportJobStatus = "expired"
)

type ExportJob struct {
	UserCtx common.UserContext

	ID              string
	TenantID        string
	RequestedBy     string
	RequestedByName string
	ResourceType    string
	ResourceLabel   string
	ProcessorKey    string
	Format          ExportJobFormat
	Status          ExportJobStatus
	ParamsJSON      string
	FileID          *string
	FileName        *string
	ErrorMessage    *string
	StartedAt       *string
	CompletedAt     *string
	ExpiresAt       *string
	CreatedAt       string
	UpdatedAt       *string
	DownloadURL     *string
}

type ExportJobListFilter struct {
	UserCtx  common.UserContext
	TenantID string
	Page     int
	Paginate int
}

type ExportJobDetailFilter struct {
	UserCtx  common.UserContext
	TenantID string
	ID       string
}

type ExportJobCreate struct {
	UserCtx common.UserContext

	TenantID      string
	RequestedBy   string
	ResourceType  string
	ResourceLabel string
	ProcessorKey  string
	Format        ExportJobFormat
	ParamsJSON    string
}

type ExportJobUpdate struct {
	UserCtx common.UserContext

	TenantID     string
	ID           string
	Status       ExportJobStatus
	FileID       *string
	ErrorMessage *string
	StartedAt    *string
	CompletedAt  *string
	ExpiresAt    *string
}

type ExportJobProcessResult struct {
	Processed int
}

type ExportJobRequestedEvent struct {
	JobID    string `json:"job_id"`
	TenantID string `json:"tenant_id"`
}
