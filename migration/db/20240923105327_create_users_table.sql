-- +goose Up
-- +goose StatementBegin
CREATE TYPE "gender" AS ENUM (
  'male',
  'female'
);

CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "role_id" uuid,
    "name" varchar(255) NOT NULL,
    "email" varchar(255) UNIQUE NOT NULL,
    "password" text,
    "phone" varchar(15),
    "gender" gender,
    "birth" timestamptz,
    "is_active" boolean NOT NULL DEFAULT true,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "created_by" uuid,
    "updated_by" uuid,
    "deleted_by" uuid,

    CONSTRAINT "fk_role_id" FOREIGN KEY ("role_id") REFERENCES "access_role" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users" DROP CONSTRAINT "fk_role_id";
DROP TABLE IF EXISTS "users";
DROP TYPE IF EXISTS "gender";
-- +goose StatementEnd