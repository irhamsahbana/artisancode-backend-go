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

func (r *exportJobRepo) GetExportJobs(ctx context.Context, filter coreentity.ExportJobListFilter) ([]coreentity.ExportJob, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:export_job:get_export_jobs:GetExportJobs")
	defer span.End()

	type dao struct {
		TotalData       int            `db:"total_data"`
		ID              string         `db:"id"`
		TenantID        string         `db:"tenant_id"`
		RequestedBy     string         `db:"requested_by"`
		RequestedByName string         `db:"requested_by_name"`
		ResourceType    string         `db:"resource_type"`
		ResourceLabel   string         `db:"resource_label"`
		ProcessorKey    string         `db:"processor_key"`
		Format          string         `db:"format"`
		Status          string         `db:"status"`
		ParamsJSON      string         `db:"params_json"`
		FileID          sql.NullString `db:"file_id"`
		ErrorMessage    sql.NullString `db:"error_message"`
		StartedAt       *time.Time     `db:"started_at"`
		CompletedAt     *time.Time     `db:"completed_at"`
		ExpiresAt       *time.Time     `db:"expires_at"`
		CreatedAt       time.Time      `db:"created_at"`
		UpdatedAt       *time.Time     `db:"updated_at"`
	}

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			are.id,
			are.tenant_id,
			are.requested_by,
			u.name AS requested_by_name,
			are.resource_type,
			are.resource_label,
			are.processor_key,
			are.format,
			are.status,
			are.params_json::text AS params_json,
			are.file_id,
			are.error_message,
			are.started_at,
			are.completed_at,
			are.expires_at,
			are.created_at,
			are.updated_at
		FROM export_jobs are
		INNER JOIN users u ON u.id = are.requested_by AND u.deleted_at IS NULL
		WHERE are.deleted_at IS NULL AND are.tenant_id = ?
		ORDER BY are.created_at DESC
		LIMIT ? OFFSET ?
	`

	args := []any{filter.TenantID, filter.Paginate, (filter.Page - 1) * filter.Paginate}
	data := make([]dao, 0)
	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query export jobs")
		return nil, 0, err
	}

	items := make([]coreentity.ExportJob, 0, len(data))
	total := 0
	for _, d := range data {
		total = d.TotalData
		item := coreentity.ExportJob{
			ID:              d.ID,
			TenantID:        d.TenantID,
			RequestedBy:     d.RequestedBy,
			RequestedByName: d.RequestedByName,
			ResourceType:    d.ResourceType,
			ResourceLabel:   d.ResourceLabel,
			ProcessorKey:    d.ProcessorKey,
			Format:          coreentity.ExportJobFormat(d.Format),
			Status:          coreentity.ExportJobStatus(d.Status),
			ParamsJSON:      d.ParamsJSON,
			ErrorMessage:    nullableStringPtr(d.ErrorMessage),
			StartedAt:       formatTimePtr(d.StartedAt),
			CompletedAt:     formatTimePtr(d.CompletedAt),
			ExpiresAt:       formatTimePtr(d.ExpiresAt),
			CreatedAt:       d.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       formatTimePtr(d.UpdatedAt),
		}
		if d.FileID.Valid {
			value := d.FileID.String
			item.FileID = &value
		}
		items = append(items, item)
	}

	return items, total, nil
}
