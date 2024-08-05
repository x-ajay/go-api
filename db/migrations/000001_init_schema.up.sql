CREATE TYPE currency AS ENUM ('USD', 'EUR', 'JPY');

CREATE TABLE IF NOT EXISTS "accounts"
(
    "id"         bigserial PRIMARY KEY,
    "owner"      TEXT        NOT NULL CHECK (length("owner") > 0),
    "balance"    INTEGER     NOT NULL,
    "currency"   currency    NOT NULL DEFAULT 'USD',
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "entries"
(
    "id"         bigserial PRIMARY KEY,
    "account_id" BIGINT      NOT NULL,
    "amount"     INTEGER     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "fk_account_id" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);

CREATE TABLE IF NOT EXISTS "transfers"
(
    "id"              bigserial PRIMARY KEY,
    "from_account_id" BIGINT      NOT NULL,
    "to_account_id"   BIGINT      NOT NULL,
    "amount"          INTEGER     NOT NULL,
    "created_at"      timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "fk_from_account_id" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id"),
    CONSTRAINT "fk_to_account_id" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id")
);

CREATE INDEX IF NOT EXISTS "idx_owner" ON "accounts" ("owner");

CREATE INDEX IF NOT EXISTS "idx_account_id" ON "entries" ("account_id");

CREATE INDEX IF NOT EXISTS "idx_from_account_id" ON "transfers" ("from_account_id");
CREATE INDEX IF NOT EXISTS "idx_to_account_id" ON "transfers" ("to_account_id");
CREATE INDEX IF NOT EXISTS "idx_from_to_account_id" ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries".amount IS 'can be negative or positive';

COMMENT ON COLUMN "transfers".amount IS 'must be positive';

