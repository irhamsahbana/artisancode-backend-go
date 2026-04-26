package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *exportJobRepo) GetExportJob(
	ctx context.Context,
	filter coreentity.ExportJobDetailFilter,
) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:export_job:get_export_job:GetExportJob",
	)
	defer span.End()

	var data struct {
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
		WHERE are.deleted_at IS NULL AND are.tenant_id = ? AND are.id = ?
		LIMIT 1
	`

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &data, exec.Rebind(query), filter.TenantID, filter.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Export job not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get export job")
		return nil, err
	}

	item := &coreentity.ExportJob{
		ID:              data.ID,
		TenantID:        data.TenantID,
		RequestedBy:     data.RequestedBy,
		RequestedByName: data.RequestedByName,
		ResourceType:    data.ResourceType,
		ResourceLabel:   data.ResourceLabel,
		ProcessorKey:    data.ProcessorKey,
		Format:          coreentity.ExportJobFormat(data.Format),
		Status:          coreentity.ExportJobStatus(data.Status),
		ParamsJSON:      data.ParamsJSON,
		ErrorMessage:    nullableStringPtr(data.ErrorMessage),
		StartedAt:       formatTimePtr(data.StartedAt),
		CompletedAt:     formatTimePtr(data.CompletedAt),
		ExpiresAt:       formatTimePtr(data.ExpiresAt),
		CreatedAt:       data.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       formatTimePtr(data.UpdatedAt),
	}
	if data.FileID.Valid {
		value := data.FileID.String
		item.FileID = &value
	}

	return item, nil
}
