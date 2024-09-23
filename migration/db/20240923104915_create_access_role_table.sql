-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "access_role" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" varchar,
    "description" varchar,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "created_by" uuid,
    "updated_by" uuid,
    "deleted_by" uuid
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "access_role";
-- +goose StatementEnd