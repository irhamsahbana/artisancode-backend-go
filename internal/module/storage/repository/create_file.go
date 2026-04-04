package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) CreateFile(ctx context.Context, data coreentity.File) (*coreentity.File, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateFile")
	defer span.End()

	query := `
		INSERT INTO storage_files (
			tenant_id,
			created_by,
			status,
			folder,
			object_key,
			original_filename,
			content_type,
			size_bytes,
			is_public
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, status, expires_at, created_at
	`

	var (
		id        string
		status    string
		expiresAt time.Time
		createdAt time.Time
	)

	err := r.db.QueryRowxContext(
		ctx,
		r.db.Rebind(query),
		data.TenantID,
		data.CreatedBy,
		common.FileStatusPending,
		data.Folder,
		data.Filename,
		data.OriginalFilename,
		data.ContentType,
		data.SizeBytes,
		data.IsPublic,
	).Scan(&id, &status, &expiresAt, &createdAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create storage file")
		return nil, err
	}

	data.ID = id
	data.Status = common.FileStatus(status)
	data.ExpiresAt = formatTimeValuePtr(expiresAt)
	data.CreatedAt = createdAt.Format(time.RFC3339)
	return &data, nil
}
