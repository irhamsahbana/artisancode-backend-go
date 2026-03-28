-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS attendance_logs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    attendance_date DATE NOT NULL,
    type VARCHAR(20) NOT NULL,
    source VARCHAR(20) NOT NULL DEFAULT 'mobile',
    status VARCHAR(20) NOT NULL DEFAULT 'recorded',
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    address TEXT,
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (tenant_id) REFERENCES tenants (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    CONSTRAINT attendance_logs_type_check CHECK (type IN ('check_in', 'check_out')),
    CONSTRAINT attendance_logs_source_check CHECK (source IN ('web', 'mobile')),
    CONSTRAINT attendance_logs_status_check CHECK (status IN ('recorded')),
    CONSTRAINT attendance_logs_unique_employee_day_type UNIQUE (tenant_id, employee_id, attendance_date, type)
);

CREATE INDEX IF NOT EXISTS idx_attendance_logs_tenant_employee_date
    ON attendance_logs (tenant_id, employee_id, attendance_date DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendance_logs;
-- +goose StatementEnd
