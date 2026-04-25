-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    role_code VARCHAR(50) NOT NULL DEFAULT 'operator',
    status VARCHAR(20) NOT NULL DEFAULT 'invited',
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT internal_users_status_check CHECK (status IN ('invited', 'active', 'inactive'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_internal_users_email_active
    ON internal_users (LOWER(email))
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_internal_users_role_status_active
    ON internal_users (role_code, status)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_users;
-- +goose StatementEnd
