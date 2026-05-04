-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_billing_reconciliation_cases (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'open',
    source_type VARCHAR(64) NOT NULL DEFAULT '',
    source_reference_id UUID,
    reason TEXT NOT NULL DEFAULT '',
    resolution_note TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    CONSTRAINT internal_billing_reconciliation_cases_status_check CHECK (status IN ('open', 'in_review', 'resolved', 'dismissed'))
);

CREATE INDEX IF NOT EXISTS idx_internal_billing_reconciliation_cases_tenant_status_active
    ON internal_billing_reconciliation_cases (tenant_id, status, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_billing_reconciliation_cases_source_active
    ON internal_billing_reconciliation_cases (source_type, source_reference_id)
    WHERE deleted_at IS NULL AND source_reference_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_billing_reconciliation_cases;
-- +goose StatementEnd
