-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "transaction" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "product_id" uuid,
  "user_id" uuid,
  "amount" int,
  "status" varchar,
  "code" varchar,
  "payment_url" varchar,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,

  CONSTRAINT "fk_product_id" FOREIGN KEY ("product_id") REFERENCES "product" ("id"),
  CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "transaction" DROP CONSTRAINT "fk_product_id";
ALTER TABLE "transaction" DROP CONSTRAINT "fk_user_id";
DROP TABLE IF EXISTS "transaction";
-- +goose StatementEnd
