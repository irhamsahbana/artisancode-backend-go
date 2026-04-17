-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_queue (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    subject VARCHAR(255) NOT NULL,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    headers_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    consumer_name VARCHAR(255),
    available_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    locked_at TIMESTAMP WITH TIME ZONE,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT message_queue_status_check CHECK (status IN ('pending', 'processing', 'processed'))
);

CREATE INDEX IF NOT EXISTS idx_message_queue_claim
    ON message_queue (status, available_at ASC, created_at ASC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_message_queue_subject_claim
    ON message_queue (subject, status, available_at ASC, created_at ASC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_message_queue_consumer_processing
    ON message_queue (consumer_name, status, locked_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_queue;
-- +goose StatementEnd
