-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_orders (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    company_id UUID,
    order_number VARCHAR(40) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending_payment',
    currency_code VARCHAR(3) NOT NULL,
    subtotal_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    total_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    internal_product_id UUID NOT NULL,
    internal_product_pricing_id UUID NOT NULL,
    pricing_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    source_type VARCHAR(32) NOT NULL,
    source_reference_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ordered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_product_id) REFERENCES internal_products (id),
    FOREIGN KEY (internal_product_pricing_id) REFERENCES internal_product_pricings (id),
    CONSTRAINT internal_orders_status_check CHECK (status IN ('draft', 'created_from_quotation', 'pending_invoice', 'pending_payment', 'paid', 'cancelled', 'expired')),
    CONSTRAINT internal_orders_source_type_check CHECK (source_type IN ('standard_pricing', 'quotation')),
    CONSTRAINT internal_orders_amount_check CHECK (
        subtotal_amount >= 0 AND discount_amount >= 0 AND tax_amount >= 0 AND total_amount >= 0
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_orders_number_active
    ON internal_orders (order_number)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_orders_tenant_status_active
    ON internal_orders (tenant_id, status)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_orders;
-- +goose StatementEnd
