-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_queue_dead_letters (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    message_queue_id UUID NOT NULL,
    subject VARCHAR(255) NOT NULL,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    headers_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    consumer_name VARCHAR(255),
    attempt_count INTEGER NOT NULL DEFAULT 0,
    failure_reason TEXT,
    queued_at TIMESTAMP WITH TIME ZONE,
    dead_lettered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_message_queue_dead_letters_message_queue_id
    ON message_queue_dead_letters (message_queue_id);

CREATE INDEX IF NOT EXISTS idx_message_queue_dead_letters_subject_dead_lettered_at
    ON message_queue_dead_letters (subject, dead_lettered_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_queue_dead_letters;
-- +goose StatementEnd
