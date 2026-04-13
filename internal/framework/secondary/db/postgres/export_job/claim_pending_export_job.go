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

func (r *exportJobRepo) ClaimPendingExportJob(ctx context.Context) (*coreentity.ExportJob, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:export_job:claim_pending_export_job:ClaimPendingExportJob")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to begin export job transaction")
		return nil, err
	}
	defer tx.Rollback()

	var data struct {
		ID            string    `db:"id"`
		TenantID      string    `db:"tenant_id"`
		RequestedBy   string    `db:"requested_by"`
		ResourceType  string    `db:"resource_type"`
		ResourceLabel string    `db:"resource_label"`
		ProcessorKey  string    `db:"processor_key"`
		Format        string    `db:"format"`
		ParamsJSON    string    `db:"params_json"`
		CreatedAt     time.Time `db:"created_at"`
	}

	query := `
		SELECT id, tenant_id, requested_by, resource_type, resource_label, processor_key, format, params_json::text AS params_json, created_at
		FROM export_jobs
		WHERE deleted_at IS NULL AND status = 'pending'
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`

	err = tx.GetContext(ctx, &data, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tx.Commit()
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to claim export job")
		return nil, err
	}

	updateQuery := `
		UPDATE export_jobs
		SET status = ?, started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP, error_message = NULL
		WHERE id = ? AND tenant_id = ?
	`

	_, err = tx.ExecContext(ctx, r.db.Rebind(updateQuery), coreentity.ExportJobStatusProcessing, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data.ID).Msg("Failed to mark export job processing")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to commit export job claim")
		return nil, err
	}

	startedAt := time.Now().UTC().Format(time.RFC3339)
	return &coreentity.ExportJob{
		ID:            data.ID,
		TenantID:      data.TenantID,
		RequestedBy:   data.RequestedBy,
		ResourceType:  data.ResourceType,
		ResourceLabel: data.ResourceLabel,
		ProcessorKey:  data.ProcessorKey,
		Format:        coreentity.ExportJobFormat(data.Format),
		Status:        coreentity.ExportJobStatusProcessing,
		ParamsJSON:    data.ParamsJSON,
		StartedAt:     &startedAt,
		CreatedAt:     data.CreatedAt.Format(time.RFC3339),
	}, nil
}
