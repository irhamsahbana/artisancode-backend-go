-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_template_roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT internal_template_roles_name_unique UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_internal_template_roles_active
    ON internal_template_roles (name)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_template_roles;
-- +goose StatementEnd
