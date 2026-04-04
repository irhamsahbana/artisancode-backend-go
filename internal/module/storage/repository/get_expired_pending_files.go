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

func (r *storageRepo) GetExpiredPendingFiles(ctx context.Context, filter coreentity.ExpiredPendingFileFilter) ([]coreentity.File, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetExpiredPendingFiles")
	defer span.End()

	before, err := time.Parse(time.RFC3339, filter.Before)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to parse expired file filter time")
		return nil, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	type dao struct {
		ID               string         `db:"id"`
		TenantID         string         `db:"tenant_id"`
		CreatedBy        string         `db:"created_by"`
		Status           string         `db:"status"`
		Folder           string         `db:"folder"`
		ObjectKey        string         `db:"object_key"`
		OriginalFilename sql.NullString `db:"original_filename"`
		ContentType      sql.NullString `db:"content_type"`
		SizeBytes        sql.NullInt64  `db:"size_bytes"`
		IsPublic         bool           `db:"is_public"`
		ExpiresAt        *time.Time     `db:"expires_at"`
		CreatedAt        time.Time      `db:"created_at"`
		UpdatedAt        *time.Time     `db:"updated_at"`
		DeletedAt        *time.Time     `db:"deleted_at"`
	}

	var data []dao
	query := `
		SELECT
			id,
			tenant_id,
			created_by,
			status,
			folder,
			object_key,
			original_filename,
			content_type,
			size_bytes,
			is_public,
			expires_at,
			created_at,
			updated_at,
			deleted_at
		FROM storage_files sf
		WHERE sf.status = ? AND sf.deleted_at IS NULL AND sf.expires_at < ?
			AND NOT EXISTS (
				SELECT 1
				FROM storage_file_links sfl
				WHERE sfl.storage_file_id = sf.id
			)
		ORDER BY sf.expires_at ASC
		LIMIT ?
	`

	err = r.db.SelectContext(ctx, &data, r.db.Rebind(query), common.FileStatusPending, before, filter.Limit)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query expired storage files")
		return nil, err
	}

	items := make([]coreentity.File, 0, len(data))
	for _, item := range data {
		file := coreentity.File{
			ID:        item.ID,
			TenantID:  item.TenantID,
			CreatedBy: item.CreatedBy,
			Status:    common.FileStatus(item.Status),
			Folder:    common.S3Folder(item.Folder),
			Filename:  item.ObjectKey,
			IsPublic:  item.IsPublic,
			ExpiresAt: formatTimePtr(item.ExpiresAt),
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: formatTimePtr(item.UpdatedAt),
			DeletedAt: formatTimePtr(item.DeletedAt),
		}
		if item.OriginalFilename.Valid {
			value := item.OriginalFilename.String
			file.OriginalFilename = &value
		}
		if item.ContentType.Valid {
			value := item.ContentType.String
			file.ContentType = &value
		}
		if item.SizeBytes.Valid {
			value := item.SizeBytes.Int64
			file.SizeBytes = &value
		}
		items = append(items, file)
	}

	return items, nil
}
