-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_payment_receipts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    internal_invoice_id UUID NOT NULL,
    internal_payment_attempt_id UUID,
    receipt_number VARCHAR(40) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending_verification',
    amount_received NUMERIC(20, 6) NOT NULL DEFAULT 0,
    currency_code VARCHAR(3) NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    verified_by_user_id UUID,
    source_type VARCHAR(64) NOT NULL DEFAULT 'manual',
    reference_number VARCHAR(128) NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (internal_invoice_id) REFERENCES internal_invoices (id),
    FOREIGN KEY (internal_payment_attempt_id) REFERENCES internal_payment_attempts (id),
    CONSTRAINT internal_payment_receipts_status_check CHECK (status IN ('pending_verification', 'accepted', 'rejected')),
    CONSTRAINT internal_payment_receipts_amount_check CHECK (amount_received >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_payment_receipts_number_active
    ON internal_payment_receipts (receipt_number)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_payment_receipts_invoice_active
    ON internal_payment_receipts (internal_invoice_id, created_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_payment_receipts;
-- +goose StatementEnd
