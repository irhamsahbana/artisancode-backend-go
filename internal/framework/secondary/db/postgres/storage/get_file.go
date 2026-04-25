package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *storageRepo) GetFile(ctx context.Context, filter coreentity.FileFilter) (*coreentity.File, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:storage:get_file:GetFile")
	defer span.End()

	var data struct {
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
		FROM storage_files
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []any{filter.TenantID}
	switch {
	case filter.ID != "":
		query += " AND id = ?"
		args = append(args, filter.ID)
	case filter.Filename != "":
		query += " AND object_key = ?"
		args = append(args, filter.Filename)
	default:
		return nil, errmsg.NewCustomErrors(400).SetMessage("File filter is required")
	}

	if filter.Folder != "" {
		query += " AND folder = ?"
		args = append(args, filter.Folder)
	}

	query = strings.TrimSpace(query) + " LIMIT 1"

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("File not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get storage file")
		return nil, err
	}

	item := coreentity.File{
		ID:        data.ID,
		TenantID:  data.TenantID,
		CreatedBy: data.CreatedBy,
		Status:    common.FileStatus(data.Status),
		Folder:    common.S3Folder(data.Folder),
		Filename:  data.ObjectKey,
		IsPublic:  data.IsPublic,
		ExpiresAt: formatTimePtr(data.ExpiresAt),
		CreatedAt: data.CreatedAt.Format(time.RFC3339),
		UpdatedAt: formatTimePtr(data.UpdatedAt),
		DeletedAt: formatTimePtr(data.DeletedAt),
	}
	if data.OriginalFilename.Valid {
		value := data.OriginalFilename.String
		item.OriginalFilename = &value
	}
	if data.ContentType.Valid {
		value := data.ContentType.String
		item.ContentType = &value
	}
	if data.SizeBytes.Valid {
		value := data.SizeBytes.Int64
		item.SizeBytes = &value
	}

	return &item, nil
}

func formatTimePtr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.Format(time.RFC3339)
	return &formatted
}

func formatTimeValuePtr(value time.Time) *string {
	formatted := value.Format(time.RFC3339)
	return &formatted
}
