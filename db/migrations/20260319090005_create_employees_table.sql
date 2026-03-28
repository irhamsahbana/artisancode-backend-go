-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    employee_no VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    user_id UUID,
    org_unit_id UUID,
    job_position_id UUID,
    location_id UUID,
    shift_id UUID,
    email VARCHAR(255),
    status VARCHAR(20) NOT NULL,
    join_date DATE,
    join_date_timezone VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (org_unit_id) REFERENCES org_units (id),
    FOREIGN KEY (job_position_id) REFERENCES job_positions (id),
    FOREIGN KEY (location_id) REFERENCES work_locations (id),
    FOREIGN KEY (shift_id) REFERENCES work_shifts (id),
    UNIQUE (tenant_id, employee_no),
    UNIQUE (tenant_id, email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd
