-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS internal_template_role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,

    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES internal_template_roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES internal_template_permissions (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS internal_template_role_permissions;
-- +goose StatementEnd
