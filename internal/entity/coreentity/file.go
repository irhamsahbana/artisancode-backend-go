package coreentity

import (
	"mime/multipart"

	"codebase-app/internal/entity/common"
)

type File struct {
	UserCtx common.UserContext

	ID        string
	TenantID  string
	Folder    common.S3Folder
	Filename  string
	URL       string
	IsPublic  bool
	CreatedAt string
}

type FileFilter struct {
	TenantID string
	Folder   common.S3Folder
	Filename string
}

type UploadFileReq struct {
	File                 *multipart.FileHeader
	Filename             string
	TenantID            string
	Folder              common.S3Folder
	IsPublic            bool
	GeneratePresignedURL bool
}

type UploadFileResp struct {
	Filename string
	URL      string
}

type DeleteFileReq struct {
	Filename string
}