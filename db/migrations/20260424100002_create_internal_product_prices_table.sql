-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_product_prices (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    internal_product_pricing_id UUID NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    amount NUMERIC(20, 6) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (internal_product_pricing_id) REFERENCES internal_product_pricings (id),
    CONSTRAINT internal_product_prices_period_check CHECK (ended_at IS NULL OR ended_at > started_at),
    CONSTRAINT internal_product_prices_amount_check CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_internal_product_prices_pricing_currency_period_active
    ON internal_product_prices (internal_product_pricing_id, currency_code, started_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_product_prices;
-- +goose StatementEnd
