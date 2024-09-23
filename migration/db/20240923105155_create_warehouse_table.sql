-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "warehouse" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" varchar,
    "address" varchar,
    "province" varchar,
    "city" varchar,
    "zip_code" varchar,
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
DROP TABLE IF EXISTS "warehouse";
-- +goose StatementEnd