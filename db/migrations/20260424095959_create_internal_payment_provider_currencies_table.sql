-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_payment_provider_currencies (
    provider VARCHAR(32) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    min_amount NUMERIC(20, 6),
    max_amount NUMERIC(20, 6),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (provider, currency_code),
    FOREIGN KEY (currency_code) REFERENCES internal_currencies (code),
    CONSTRAINT internal_payment_provider_currencies_provider_check CHECK (provider ~ '^[a-z0-9_-]+$'),
    CONSTRAINT internal_payment_provider_currencies_amount_check CHECK (
        (min_amount IS NULL OR min_amount >= 0)
        AND (max_amount IS NULL OR max_amount >= 0)
        AND (min_amount IS NULL OR max_amount IS NULL OR max_amount >= min_amount)
    )
);

CREATE INDEX IF NOT EXISTS idx_internal_payment_provider_currencies_active
    ON internal_payment_provider_currencies (provider, is_active, currency_code)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_payment_provider_currencies;
-- +goose StatementEnd
