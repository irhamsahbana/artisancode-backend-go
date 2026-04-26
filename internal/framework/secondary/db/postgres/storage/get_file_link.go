package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) GetFileLink(
	ctx context.Context,
	filter coreentity.StorageFileLinkFilter,
) (*coreentity.StorageFileLink, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:storage:get_file_link:GetFileLink")
	defer span.End()

	var data struct {
		ID            string    `db:"id"`
		TenantID      string    `db:"tenant_id"`
		StorageFileID string    `db:"storage_file_id"`
		ResourceType  string    `db:"resource_type"`
		ResourceID    string    `db:"resource_id"`
		FieldName     string    `db:"field_name"`
		SortOrder     int       `db:"sort_order"`
		CreatedAt     time.Time `db:"created_at"`
	}

	query := `
		SELECT
			id,
			tenant_id,
			storage_file_id,
			resource_type,
			resource_id,
			field_name,
			sort_order,
			created_at
		FROM storage_file_links
		WHERE tenant_id = ?
			AND resource_type = ?
			AND resource_id = ?
			AND field_name = ?
		ORDER BY sort_order ASC, created_at ASC
		LIMIT 1
	`

	exec := r.executor(ctx)
	err := exec.GetContext(
		ctx,
		&data,
		exec.Rebind(query),
		filter.TenantID,
		filter.ResourceType,
		filter.ResourceID,
		filter.FieldName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get storage file link")
		return nil, err
	}

	return &coreentity.StorageFileLink{
		ID:            data.ID,
		TenantID:      data.TenantID,
		StorageFileID: data.StorageFileID,
		ResourceType:  data.ResourceType,
		ResourceID:    data.ResourceID,
		FieldName:     data.FieldName,
		SortOrder:     data.SortOrder,
		CreatedAt:     data.CreatedAt.Format(time.RFC3339),
	}, nil
}
