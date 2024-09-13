-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "gender" AS ENUM (
  'pria',
  'wanita'
);

CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" varchar(255) NOT NULL,
    "email" varchar(255) UNIQUE NOT NULL,
    "occupation" varchar,
    "password" text,
    "phone" varchar(15),
    "gender" gender,
    "role" varchar,
    "token" text NOT NULL,
    "is_google" boolean DEFAULT false,
    "is_active" boolean NOT NULL DEFAULT true,
    "is_verify" boolean DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
DROP TYPE IF EXISTS "gender";
-- +goose StatementEnd
