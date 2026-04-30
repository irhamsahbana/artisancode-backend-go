-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_auth_identities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    provider VARCHAR(32) NOT NULL,
    provider_subject VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT false,
    display_name VARCHAR(255) NOT NULL DEFAULT '',
    picture_url TEXT,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_identities_provider_subject_active
    ON user_auth_identities (provider, provider_subject)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_identities_user_provider_active
    ON user_auth_identities (tenant_id, user_id, provider)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_code_active_unique
    ON tenants (code)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tenants_code_active_unique;
DROP INDEX IF EXISTS idx_user_auth_identities_user_provider_active;
DROP INDEX IF EXISTS idx_user_auth_identities_provider_subject_active;
DROP TABLE IF EXISTS user_auth_identities;
-- +goose StatementEnd
