-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS export_jobs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    requested_by UUID NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_label VARCHAR(255) NOT NULL,
    processor_key VARCHAR(100) NOT NULL,
    format VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    params_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    file_id UUID,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (requested_by) REFERENCES users (id),
    FOREIGN KEY (file_id) REFERENCES storage_files (id),
    CONSTRAINT export_jobs_format_check CHECK (format IN ('csv', 'xlsx', 'pdf')),
    CONSTRAINT export_jobs_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_export_jobs_tenant_created
    ON export_jobs (tenant_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_export_jobs_pending
    ON export_jobs (status, created_at ASC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS export_jobs;
-- +goose StatementEnd
