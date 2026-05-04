-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_billing_accounts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'open',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    CONSTRAINT internal_billing_accounts_status_check CHECK (status IN ('open', 'suspended', 'closed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_billing_accounts_tenant_active
    ON internal_billing_accounts (tenant_id)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_billing_accounts;
-- +goose StatementEnd
