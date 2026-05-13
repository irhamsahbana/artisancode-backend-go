-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    org_unit_id UUID,
    parent_id UUID,
    group_key VARCHAR(100) NOT NULL, -- logical grouping across reusable category sets, e.g. expense_category or leave_reason
    code VARCHAR(100) NOT NULL, -- stable internal identifier; safe for seeds, filters, and logic even when name changes
    name VARCHAR(255) NOT NULL,
    notes TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (org_unit_id) REFERENCES org_units (id),
    FOREIGN KEY (parent_id) REFERENCES categories (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_tenant_level_unique
    ON categories (tenant_id, group_key, code)
    WHERE deleted_at IS NULL AND org_unit_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_org_unit_level_unique
    ON categories (tenant_id, org_unit_id, group_key, code)
    WHERE deleted_at IS NULL AND org_unit_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_categories_tenant_id_active
    ON categories (tenant_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_categories_tenant_group_key_parent_active
    ON categories (tenant_id, group_key, parent_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_categories_tenant_org_unit_group_key_active
    ON categories (tenant_id, org_unit_id, group_key)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
