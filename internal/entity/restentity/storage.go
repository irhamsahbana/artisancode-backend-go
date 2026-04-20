package restentity

import (
	"codebase-app/internal/entity/common"
	"mime/multipart"
)

type UploadFileReq struct {
	File                 *multipart.FileHeader `form:"file" validate:"required"`
	Filename             string
	Folder               common.S3Folder
	IsPublic             bool
	GeneratePresignedURL bool
}

type GetFileReq struct {
	FileName string `json:"filename" validate:"required"`
}

type UploadFileResp struct {
	FileName string `json:"filename" validate:"required"`
	Url      string `json:"url" validate:"required"`
}

type DeleteFileReq struct {
	FileName string `json:"filename" validate:"required"`
}

type CreateUploadURLReq struct {
	Filename         string          `json:"filename" validate:"required,max=255"`
	OriginalFilename *string         `json:"original_filename" validate:"omitempty,max=255"`
	ContentType      string          `json:"content_type" validate:"required,max=255"`
	Folder           common.S3Folder `json:"folder" validate:"required"`
	IsPublic         bool            `json:"is_public"`
}

type CreateUploadURLResp struct {
	FileID    string            `json:"file_id"`
	ObjectKey string            `json:"object_key"`
	UploadURL string            `json:"upload_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
}
