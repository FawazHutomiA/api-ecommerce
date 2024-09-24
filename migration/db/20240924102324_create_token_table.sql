-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "token" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" uuid,
    "token" varchar,
    "expired_at" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "created_by" uuid,
    "updated_by" uuid,
    "deleted_by" uuid,

    CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "token" DROP CONSTRAINT "fk_user_id";
DROP TABLE IF EXISTS "token";
-- +goose StatementEnd