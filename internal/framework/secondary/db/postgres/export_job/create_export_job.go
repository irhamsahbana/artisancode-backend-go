package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *exportJobRepo) CreateExportJob(
	ctx context.Context,
	data coreentity.ExportJobCreate,
) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:export_job:create_export_job:CreateExportJob",
	)
	defer span.End()

	query := `
		INSERT INTO export_jobs (
			tenant_id,
			requested_by,
			resource_type,
			resource_label,
			processor_key,
			format,
			status,
			params_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?::jsonb)
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
		data.TenantID,
		data.RequestedBy,
		data.ResourceType,
		data.ResourceLabel,
		data.ProcessorKey,
		data.Format,
		coreentity.ExportJobStatusPending,
		data.ParamsJSON,
	).Scan(&id, &createdAt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create export job")
		return nil, err
	}

	return &coreentity.ExportJob{
		ID:            id,
		TenantID:      data.TenantID,
		RequestedBy:   data.RequestedBy,
		ResourceType:  data.ResourceType,
		ResourceLabel: data.ResourceLabel,
		ProcessorKey:  data.ProcessorKey,
		Format:        data.Format,
		Status:        coreentity.ExportJobStatusPending,
		ParamsJSON:    data.ParamsJSON,
		CreatedAt:     createdAt.Format(time.RFC3339),
	}, nil
}
