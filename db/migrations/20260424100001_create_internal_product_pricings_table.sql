-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_product_pricings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    internal_product_id UUID NOT NULL,
    code VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (internal_product_id) REFERENCES internal_products (id),
    CONSTRAINT internal_product_pricings_status_check CHECK (status IN ('draft', 'active', 'inactive', 'archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_product_pricings_product_code_active
    ON internal_product_pricings (internal_product_id, code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_product_pricings_product_name_active
    ON internal_product_pricings (internal_product_id, name)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_product_pricings;
-- +goose StatementEnd
