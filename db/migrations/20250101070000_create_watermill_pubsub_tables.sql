-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "watermill_email_verification" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_email_verification" (
    "consumer_group" VARCHAR(255) NOT NULL,
    "offset_acked" BIGINT,
    "last_processed_transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("consumer_group")
);

CREATE TABLE IF NOT EXISTS "watermill_email_forgot_password" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_email_forgot_password" (
    "consumer_group" VARCHAR(255) NOT NULL,
    "offset_acked" BIGINT,
    "last_processed_transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("consumer_group")
);

CREATE TABLE IF NOT EXISTS "watermill_email_invitation" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_email_invitation" (
    "consumer_group" VARCHAR(255) NOT NULL,
    "offset_acked" BIGINT,
    "last_processed_transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("consumer_group")
);

CREATE TABLE IF NOT EXISTS "watermill_export_job_requested" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_export_job_requested" (
    "consumer_group" VARCHAR(255) NOT NULL,
    "offset_acked" BIGINT,
    "last_processed_transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("consumer_group")
);

CREATE TABLE IF NOT EXISTS "watermill_dead_letter" (
    "offset" BIGSERIAL,
    "uuid" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "payload" JSON DEFAULT NULL,
    "metadata" JSON DEFAULT NULL,
    "transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("transaction_id", "offset")
);

CREATE TABLE IF NOT EXISTS "watermill_offsets_dead_letter" (
    "consumer_group" VARCHAR(255) NOT NULL,
    "offset_acked" BIGINT,
    "last_processed_transaction_id" xid8 NOT NULL,
    PRIMARY KEY ("consumer_group")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "watermill_offsets_dead_letter";
DROP TABLE IF EXISTS "watermill_dead_letter";
DROP TABLE IF EXISTS "watermill_offsets_export_job_requested";
DROP TABLE IF EXISTS "watermill_export_job_requested";
DROP TABLE IF EXISTS "watermill_offsets_email_invitation";
DROP TABLE IF EXISTS "watermill_email_invitation";
DROP TABLE IF EXISTS "watermill_offsets_email_forgot_password";
DROP TABLE IF EXISTS "watermill_email_forgot_password";
DROP TABLE IF EXISTS "watermill_offsets_email_verification";
DROP TABLE IF EXISTS "watermill_email_verification";
-- +goose StatementEnd
