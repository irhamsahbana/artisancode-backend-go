-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_tenant_invoices (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_billing_account_id UUID NOT NULL,
    internal_tenant_subscription_id UUID,
    internal_tenant_subscription_change_id UUID,
    invoice_number VARCHAR(40) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'open',
    currency_code VARCHAR(3) NOT NULL,
    amount NUMERIC(20, 6) NOT NULL DEFAULT 0,
    amount_paid NUMERIC(20, 6) NOT NULL DEFAULT 0,
    amount_outstanding NUMERIC(20, 6) NOT NULL DEFAULT 0,
    source_type VARCHAR(64) NOT NULL DEFAULT 'self_serve_checkout',
    source_reference_id UUID,
    target_subscription_state VARCHAR(32) NOT NULL DEFAULT 'pending_activation',
    due_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_billing_account_id) REFERENCES internal_billing_accounts (id),
    FOREIGN KEY (internal_tenant_subscription_id) REFERENCES internal_tenant_subscriptions (id),
    FOREIGN KEY (internal_tenant_subscription_change_id) REFERENCES internal_tenant_subscription_changes (id),
    CONSTRAINT internal_tenant_invoices_status_check CHECK (status IN ('draft', 'open', 'partially_paid', 'paid', 'expired', 'cancelled')),
    CONSTRAINT internal_tenant_invoices_target_subscription_state_check CHECK (target_subscription_state IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled', 'expired')),
    CONSTRAINT internal_tenant_invoices_source_type_check CHECK (source_type IN ('self_serve_checkout', 'assisted_manual', 'renewal_scheduler')),
    CONSTRAINT internal_tenant_invoices_amount_check CHECK (
        amount >= 0
        AND amount_paid >= 0
        AND amount_outstanding >= 0
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_tenant_invoices_number_active
    ON internal_tenant_invoices (invoice_number)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_tenant_invoices_tenant_status_active
    ON internal_tenant_invoices (tenant_id, status, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_tenant_invoices_change_active
    ON internal_tenant_invoices (internal_tenant_subscription_change_id)
    WHERE deleted_at IS NULL AND internal_tenant_subscription_change_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_tenant_invoices;
-- +goose StatementEnd
