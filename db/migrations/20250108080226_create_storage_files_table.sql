-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS storage_files (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    created_by UUID NOT NULL,
    created_by_name VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    folder VARCHAR(100) NOT NULL,
    object_key TEXT NOT NULL,
    original_filename VARCHAR(255),
    content_type VARCHAR(255),
    size_bytes BIGINT,
    is_public BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT storage_files_status_check CHECK (status IN ('pending', 'attached', 'deleted', 'failed')),
    UNIQUE (tenant_id, object_key)
);

CREATE INDEX IF NOT EXISTS idx_storage_files_status_expires_at
    ON storage_files (status, expires_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_storage_files_tenant_id_active
    ON storage_files (tenant_id)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS storage_files;
-- +goose StatementEnd
