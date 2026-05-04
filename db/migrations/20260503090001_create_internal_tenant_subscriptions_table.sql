-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_tenant_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    internal_billing_account_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    product_family VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'free',
    internal_product_id UUID,
    internal_product_pricing_id UUID,
    plan_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    add_on_snapshots JSONB NOT NULL DEFAULT '[]'::jsonb,
    current_period_started_at TIMESTAMP WITH TIME ZONE,
    current_period_ended_at TIMESTAMP WITH TIME ZONE,
    grace_ended_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (internal_billing_account_id) REFERENCES internal_billing_accounts (id),
    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_product_id) REFERENCES internal_products (id),
    FOREIGN KEY (internal_product_pricing_id) REFERENCES internal_product_pricings (id),
    CONSTRAINT internal_tenant_subscriptions_status_check CHECK (status IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled', 'expired')),
    CONSTRAINT internal_tenant_subscriptions_period_check CHECK (
        current_period_ended_at IS NULL
        OR current_period_started_at IS NULL
        OR current_period_ended_at > current_period_started_at
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_tenant_subscriptions_tenant_family_current
    ON internal_tenant_subscriptions (tenant_id, product_family)
    WHERE deleted_at IS NULL AND status IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled');

CREATE INDEX IF NOT EXISTS idx_internal_tenant_subscriptions_account_active
    ON internal_tenant_subscriptions (internal_billing_account_id, created_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_tenant_subscriptions;
-- +goose StatementEnd
