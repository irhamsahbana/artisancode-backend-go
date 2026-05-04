-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_billing_ledger_entries (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_tenant_subscription_id UUID,
    entry_type VARCHAR(40) NOT NULL,
    source_type VARCHAR(64) NOT NULL DEFAULT '',
    source_reference_id UUID,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_tenant_subscription_id) REFERENCES internal_tenant_subscriptions (id),
    CONSTRAINT internal_billing_ledger_entries_type_check CHECK (entry_type IN ('invoice_issued', 'payment_succeeded', 'payment_failed', 'subscription_change', 'entitlement_updated', 'adjustment'))
);

CREATE INDEX IF NOT EXISTS idx_internal_billing_ledger_entries_tenant_occurred_active
    ON internal_billing_ledger_entries (tenant_id, occurred_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_billing_ledger_entries_source_active
    ON internal_billing_ledger_entries (source_type, source_reference_id)
    WHERE deleted_at IS NULL AND source_reference_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_billing_ledger_entries;
-- +goose StatementEnd
