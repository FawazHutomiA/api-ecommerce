-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "product" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" uuid,
  "name" varchar,
  "short_description" varchar,
  "description" varchar,
  "goal_amount" int,
  "current_amount" int,
  "slug" varchar,
  "backer_amount" int,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,

  CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "product" DROP CONSTRAINT "fk_user_id";
DROP TABLE IF EXISTS "product";
-- +goose StatementEnd
