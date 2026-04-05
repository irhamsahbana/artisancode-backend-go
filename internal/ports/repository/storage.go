package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type StorageRepository interface {
	CreateFile(ctx context.Context, data coreentity.File) (*coreentity.File, error)
	GetFile(ctx context.Context, filter coreentity.FileFilter) (*coreentity.File, error)
	GetFileLink(ctx context.Context, filter coreentity.StorageFileLinkFilter) (*coreentity.StorageFileLink, error)
	CreateFileLink(ctx context.Context, req coreentity.CreateStorageFileLinkReq) (*coreentity.StorageFileLink, error)
	MarkFileAttached(ctx context.Context, tenantID, id string) error
	GetExpiredPendingFiles(ctx context.Context, filter coreentity.ExpiredPendingFileFilter) ([]coreentity.File, error)
	MarkFileDeleted(ctx context.Context, tenantID, id string) error
}
