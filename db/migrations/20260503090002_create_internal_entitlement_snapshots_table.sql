-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_entitlement_snapshots (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_tenant_subscription_id UUID NOT NULL,
    subscription_status VARCHAR(32) NOT NULL,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    usage_limits JSONB NOT NULL DEFAULT '{}'::jsonb,
    source_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    effective_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_tenant_subscription_id) REFERENCES internal_tenant_subscriptions (id),
    CONSTRAINT internal_entitlement_snapshots_status_check CHECK (subscription_status IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_internal_entitlement_snapshots_tenant_effective_active
    ON internal_entitlement_snapshots (tenant_id, effective_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_entitlement_snapshots_subscription_effective_active
    ON internal_entitlement_snapshots (internal_tenant_subscription_id, effective_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_entitlement_snapshots;
-- +goose StatementEnd
