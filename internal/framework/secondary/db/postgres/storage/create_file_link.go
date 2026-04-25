package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) CreateFileLink(ctx context.Context, req coreentity.CreateStorageFileLinkReq) (*coreentity.StorageFileLink, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:storage:create_file_link:CreateFileLink")
	defer span.End()

	query := `
		INSERT INTO storage_file_links (
			tenant_id,
			storage_file_id,
			resource_type,
			resource_id,
			field_name,
			sort_order
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	var (
		id        string
		createdAt time.Time
	)

	exec := r.executor(ctx)
	err := exec.QueryRowxContext(
		ctx,
		exec.Rebind(query),
		req.TenantID,
		req.StorageFileID,
		req.ResourceType,
		req.ResourceID,
		req.FieldName,
		req.SortOrder,
	).Scan(&id, &createdAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("Failed to create storage file link")
		return nil, err
	}

	return &coreentity.StorageFileLink{
		ID:            id,
		TenantID:      req.TenantID,
		StorageFileID: req.StorageFileID,
		ResourceType:  req.ResourceType,
		ResourceID:    req.ResourceID,
		FieldName:     req.FieldName,
		SortOrder:     req.SortOrder,
		CreatedAt:     createdAt.Format(time.RFC3339),
	}, nil
}
