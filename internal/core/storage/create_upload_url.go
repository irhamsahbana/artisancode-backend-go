package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *storageCore) CreateUploadURL(
	ctx context.Context,
	req coreentity.PresignUploadURLReq,
) (*coreentity.PresignUploadURLResp, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:storage:create_upload_url:CreateUploadURL")
	defer span.End()

	objectKey := buildObjectKey(req.TenantID, req.CreatedBy, req.Folder, req.Filename, req.ContentType, req.IsPublic)
	req.Filename = objectKey

	file, err := c.repo.CreateFile(ctx, coreentity.File{
		UserCtx:          req.UserCtx,
		TenantID:         req.TenantID,
		CreatedBy:        req.CreatedBy,
		Folder:           req.Folder,
		Filename:         req.Filename,
		OriginalFilename: req.OriginalFilename,
		ContentType:      &req.ContentType,
		IsPublic:         req.IsPublic,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.s3.PresignUploadURL(ctx, &req)
	if err != nil {
		return nil, err
	}
	resp.FileID = file.ID
	return resp, nil
}

func buildObjectKey(
	tenantID, userID string,
	folder common.S3Folder,
	filename, contentType string,
	isPublic bool,
) string {
	extension := ".jpg"
	switch strings.ToLower(contentType) {
	case "image/png":
		extension = ".png"
	case "image/heic":
		extension = ".heic"
	}

	visibilityFolder := "private"
	if isPublic {
		visibilityFolder = "public"
	}

	safeName := strings.TrimSuffix(filename, extension)
	safeName = strings.ReplaceAll(safeName, " ", "-")

	return fmt.Sprintf(
		"%s/tenants/%s/%s/users/%s/%d-%s%s",
		visibilityFolder,
		tenantID,
		strings.TrimPrefix(string(folder), "/"),
		userID,
		time.Now().UTC().UnixNano(),
		safeName,
		extension,
	)
}
