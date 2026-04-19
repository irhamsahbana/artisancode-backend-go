-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS work_locations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    org_unit_id UUID,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10,6),
    longitude DECIMAL(10,6),
    radius_meters INT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (org_unit_id) REFERENCES org_units (id),
    CONSTRAINT work_locations_name_unique UNIQUE (tenant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_work_locations_tenant_id_active
    ON work_locations (tenant_id)
    WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS work_locations;
-- +goose StatementEnd
