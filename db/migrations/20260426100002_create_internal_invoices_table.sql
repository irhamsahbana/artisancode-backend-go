-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_invoices (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    internal_order_id UUID NOT NULL,
    invoice_number VARCHAR(40) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'open',
    currency_code VARCHAR(3) NOT NULL,
    amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    amount_paid NUMERIC(20, 6) NOT NULL DEFAULT 0,
    amount_outstanding NUMERIC(20, 6) NOT NULL DEFAULT 0,
    due_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (internal_order_id) REFERENCES internal_orders (id),
    CONSTRAINT internal_invoices_status_check CHECK (status IN ('draft', 'open', 'partially_paid', 'paid', 'expired', 'cancelled')),
    CONSTRAINT internal_invoices_amount_check CHECK (amount >= 0 AND amount_paid >= 0 AND amount_outstanding >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_invoices_number_active
    ON internal_invoices (invoice_number)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_invoices_order_active
    ON internal_invoices (internal_order_id)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_invoices;
-- +goose StatementEnd
