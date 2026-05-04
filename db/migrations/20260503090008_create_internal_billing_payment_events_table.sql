-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_billing_payment_events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_tenant_invoice_id UUID,
    internal_tenant_payment_attempt_id UUID,
    provider VARCHAR(32) NOT NULL,
    provider_event_id VARCHAR(160) NOT NULL,
    provider_request_id VARCHAR(128) NOT NULL DEFAULT '',
    provider_reference VARCHAR(128) NOT NULL DEFAULT '',
    provider_status VARCHAR(128) NOT NULL,
    provider_context VARCHAR(64) NOT NULL,
    normalized_attempt_status VARCHAR(24) NOT NULL DEFAULT '',
    normalized_event_status VARCHAR(64) NOT NULL DEFAULT '',
    is_terminal BOOLEAN NOT NULL DEFAULT FALSE,
    processing_result VARCHAR(64) NOT NULL DEFAULT 'received',
    processed_at TIMESTAMP WITH TIME ZONE,
    raw_payload_summary JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_tenant_invoice_id) REFERENCES internal_tenant_invoices (id),
    FOREIGN KEY (internal_tenant_payment_attempt_id) REFERENCES internal_tenant_payment_attempts (id),
    CONSTRAINT internal_billing_payment_events_attempt_status_check CHECK (
        normalized_attempt_status = ''
        OR normalized_attempt_status IN ('initiated', 'pending', 'succeeded', 'failed', 'expired', 'cancelled')
    ),
    CONSTRAINT internal_billing_payment_events_event_status_check CHECK (
        normalized_event_status = ''
        OR normalized_event_status IN ('failed_by_channel', 'recovered', 'cancelled_by_merchant')
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_billing_payment_events_provider_event_active
    ON internal_billing_payment_events (provider, provider_event_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_billing_payment_events_attempt_active
    ON internal_billing_payment_events (internal_tenant_payment_attempt_id, created_at DESC)
    WHERE deleted_at IS NULL AND internal_tenant_payment_attempt_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_internal_billing_payment_events_tenant_active
    ON internal_billing_payment_events (tenant_id, created_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_billing_payment_events;
-- +goose StatementEnd
