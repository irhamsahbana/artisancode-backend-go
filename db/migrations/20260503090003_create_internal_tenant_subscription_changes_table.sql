-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_tenant_subscription_changes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    internal_tenant_subscription_id UUID NOT NULL,
    change_type VARCHAR(32) NOT NULL,
    from_status VARCHAR(32) NOT NULL,
    to_status VARCHAR(32) NOT NULL,
    trigger VARCHAR(64) NOT NULL,
    source_type VARCHAR(64) NOT NULL DEFAULT '',
    source_reference_id UUID,
    effective_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (internal_tenant_subscription_id) REFERENCES internal_tenant_subscriptions (id),
    CONSTRAINT internal_tenant_subscription_changes_type_check CHECK (change_type IN ('checkout', 'renewal', 'upgrade', 'downgrade', 'cancellation', 'reactivation', 'expiration')),
    CONSTRAINT internal_tenant_subscription_changes_from_status_check CHECK (from_status IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled', 'expired')),
    CONSTRAINT internal_tenant_subscription_changes_to_status_check CHECK (to_status IN ('free', 'pending_activation', 'active', 'grace_period', 'suspended', 'cancelled', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_internal_tenant_subscription_changes_subscription_effective_active
    ON internal_tenant_subscription_changes (internal_tenant_subscription_id, effective_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_tenant_subscription_changes_source_active
    ON internal_tenant_subscription_changes (source_type, source_reference_id)
    WHERE deleted_at IS NULL AND source_reference_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_tenant_subscription_changes;
-- +goose StatementEnd
