-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    symbol VARCHAR(16) NOT NULL DEFAULT '',
    decimal_places SMALLINT NOT NULL DEFAULT 2,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT internal_currencies_code_check CHECK (code ~ '^[A-Z]{3}$'),
    CONSTRAINT internal_currencies_decimal_places_check CHECK (decimal_places >= 0 AND decimal_places <= 6)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_currencies_default_active
    ON internal_currencies (is_default)
    WHERE is_default = TRUE AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_currencies_active_sort
    ON internal_currencies (is_active, sort_order, code)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_currencies;
-- +goose StatementEnd
