-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS storage_file_links (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    storage_file_id UUID NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id TEXT NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    sort_order INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (storage_file_id) REFERENCES storage_files (id),
    UNIQUE (storage_file_id, resource_type, resource_id, field_name, sort_order)
);

CREATE INDEX IF NOT EXISTS idx_storage_file_links_resource
    ON storage_file_links (tenant_id, resource_type, resource_id, field_name, sort_order);

CREATE INDEX IF NOT EXISTS idx_storage_file_links_storage_file_id
    ON storage_file_links (storage_file_id);

CREATE INDEX IF NOT EXISTS idx_storage_file_links_tenant_id_storage_file_id
    ON storage_file_links (tenant_id, storage_file_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS storage_file_links;
-- +goose StatementEnd
