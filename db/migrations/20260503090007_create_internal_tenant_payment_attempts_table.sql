-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_tenant_payment_attempts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_tenant_invoice_id UUID NOT NULL,
    provider VARCHAR(32) NOT NULL,
    payment_method_type VARCHAR(64) NOT NULL,
    payment_channel_code VARCHAR(64) NOT NULL DEFAULT '',
    provider_reference VARCHAR(128) NOT NULL DEFAULT '',
    provider_request_id VARCHAR(128) NOT NULL DEFAULT '',
    provider_payment_url TEXT NOT NULL DEFAULT '',
    provider_payload_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(24) NOT NULL DEFAULT 'initiated',
    requested_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    expired_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    raw_last_status VARCHAR(128) NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_tenant_invoice_id) REFERENCES internal_tenant_invoices (id),
    CONSTRAINT internal_tenant_payment_attempts_status_check CHECK (status IN ('initiated', 'pending', 'succeeded', 'failed', 'expired', 'cancelled')),
    CONSTRAINT internal_tenant_payment_attempts_amount_check CHECK (
        requested_amount >= 0
        AND paid_amount >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_internal_tenant_payment_attempts_invoice_active
    ON internal_tenant_payment_attempts (internal_tenant_invoice_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_tenant_payment_attempts_provider_request_active
    ON internal_tenant_payment_attempts (provider, provider_request_id)
    WHERE deleted_at IS NULL AND provider_request_id <> '';

CREATE INDEX IF NOT EXISTS idx_internal_tenant_payment_attempts_provider_reference_active
    ON internal_tenant_payment_attempts (provider, provider_reference)
    WHERE deleted_at IS NULL AND provider_reference <> '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_tenant_payment_attempts;
-- +goose StatementEnd
