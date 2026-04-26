-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_quotations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    company_id UUID,
    quotation_number VARCHAR(40) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'draft',
    currency_code VARCHAR(3) NOT NULL,
    subtotal_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    total_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    internal_product_id UUID NOT NULL,
    internal_product_pricing_id UUID NOT NULL,
    pricing_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    quote_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    expires_at TIMESTAMP WITH TIME ZONE,
    approved_at TIMESTAMP WITH TIME ZONE,
    converted_to_order_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_product_id) REFERENCES internal_products (id),
    FOREIGN KEY (internal_product_pricing_id) REFERENCES internal_product_pricings (id),
    CONSTRAINT internal_quotations_status_check CHECK (status IN ('draft', 'sent', 'approved', 'rejected', 'expired', 'converted')),
    CONSTRAINT internal_quotations_amount_check CHECK (
        subtotal_amount >= 0 AND discount_amount >= 0 AND tax_amount >= 0 AND total_amount >= 0
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_quotations_number_active
    ON internal_quotations (quotation_number)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_quotations_tenant_status_active
    ON internal_quotations (tenant_id, status)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_quotations;
-- +goose StatementEnd
