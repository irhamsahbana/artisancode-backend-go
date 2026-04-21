-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_invitations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    employee_id UUID,
    email VARCHAR(255) NOT NULL,
    role_code VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    invited_by UUID NOT NULL,
    last_sent_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    FOREIGN KEY (invited_by) REFERENCES users (id),
    CONSTRAINT user_invitations_status_check CHECK (status IN ('pending', 'accepted', 'expired', 'revoked'))
);

CREATE INDEX IF NOT EXISTS idx_user_invitations_tenant_id_active
    ON user_invitations (tenant_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_invitations_token_hash_active
    ON user_invitations (token_hash)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_invitations_employee_id_active
    ON user_invitations (employee_id)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_invitations;
-- +goose StatementEnd
