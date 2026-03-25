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
