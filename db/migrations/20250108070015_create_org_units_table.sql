-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS org_units (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID,
    parent_id UUID,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(64) NOT NULL,
    category VARCHAR(50) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT org_units_code_unique UNIQUE (tenant_id, code),
    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (parent_id) REFERENCES org_units (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS org_units;
-- +goose StatementEnd
