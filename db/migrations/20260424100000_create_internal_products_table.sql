-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_products (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    code VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT internal_products_status_check CHECK (status IN ('draft', 'active', 'inactive', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_products_code_active
    ON internal_products (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_products_name_active
    ON internal_products (name)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_products;
-- +goose StatementEnd
