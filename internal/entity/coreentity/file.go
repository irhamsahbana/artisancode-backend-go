package coreentity

import (
	"mime/multipart"

	"codebase-app/internal/entity/common"
)

type File struct {
	UserCtx common.UserContext

	ID               string
	TenantID         string
	CreatedBy        string
	Status           common.FileStatus
	Folder           common.S3Folder
	Filename         string
	OriginalFilename *string
	ContentType      *string
	SizeBytes        *int64
	URL              string
	IsPublic         bool
	ExpiresAt        *string
	CreatedAt        string
	UpdatedAt        *string
	DeletedAt        *string
}

type FileFilter struct {
	TenantID string
	ID       string
	Folder   common.S3Folder
	Filename string
}

type UploadFileReq struct {
	File                 *multipart.FileHeader
	Filename             string
	TenantID             string
	Folder               common.S3Folder
	IsPublic             bool
	GeneratePresignedURL bool
}

type UploadFileResp struct {
	Filename string
	URL      string
}

type PresignUploadURLReq struct {
	UserCtx          common.UserContext
	TenantID         string
	CreatedBy        string
	Filename         string
	OriginalFilename *string
	ContentType      string
	Folder           common.S3Folder
	IsPublic         bool
}

type ExpiredPendingFileFilter struct {
	Before string
	Limit  int
}

type CleanupExpiredFilesReq struct {
	Before string
	Limit  int
}

type CleanupExpiredFilesResp struct {
	Scanned int
	Deleted int
}

type StorageFileLink struct {
	ID            string
	TenantID      string
	StorageFileID string
	ResourceType  string
	ResourceID    string
	FieldName     string
	SortOrder     int
	CreatedAt     string
}

type StorageFileLinkFilter struct {
	TenantID     string
	ResourceType string
	ResourceID   string
	FieldName    string
}

type CreateStorageFileLinkReq struct {
	TenantID      string
	StorageFileID string
	ResourceType  string
	ResourceID    string
	FieldName     string
	SortOrder     int
}

type PresignUploadURLResp struct {
	FileID   string
	Filename string
	URL      string
	Method   string
	Headers  map[string]string
}

type DeleteFileReq struct {
	TenantID string
	Filename string
}
