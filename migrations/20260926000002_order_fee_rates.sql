-- +goose Up
-- ============================================================
-- 1) fee_rates: time-windowed shipping/packaging tariffs managed
--    from the admin panel (هزینه ارسال / هزینه بسته‌بندی).
--    A new row per period keeps history; the active row is resolved
--    by status + [starts_at, ends_at] at checkout time.
-- 2) orders: fee breakdown snapshot + grand total. Existing rows are
--    backfilled so GrandTotal always equals what was (or will be) paid.
-- ============================================================

CREATE TABLE fee_rates (
    id             BIGSERIAL PRIMARY KEY,
    kind           VARCHAR(16) NOT NULL CHECK (kind IN ('shipping', 'packaging')),
    title          TEXT NOT NULL,
    amount         BIGINT NOT NULL DEFAULT 0,
    -- shipping only: orders with an items total >= threshold ship free.
    -- NULL/0 disables the rule.
    free_threshold BIGINT,
    starts_at      DATE,
    ends_at        DATE,
    status         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT now(),
    updated_at     TIMESTAMPTZ DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);
CREATE INDEX idx_fee_rates_active ON fee_rates (kind, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_fee_rates_deleted_at ON fee_rates (deleted_at);

COMMENT ON COLUMN fee_rates.kind IS 'shipping | packaging';
COMMENT ON COLUMN fee_rates.free_threshold IS
    'shipping only: items total at/above which shipping becomes free';

ALTER TABLE orders ADD COLUMN shipping_fee   BIGINT NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN packaging_fee  BIGINT NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN shipping_free  BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE orders ADD COLUMN grand_total    BIGINT NOT NULL DEFAULT 0;

-- existing orders were paid for items only → grand total equals that
UPDATE orders SET grand_total = total_sale_price WHERE grand_total = 0;

-- +goose Down
ALTER TABLE orders DROP COLUMN IF EXISTS grand_total;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_free;
ALTER TABLE orders DROP COLUMN IF EXISTS packaging_fee;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_fee;

DROP TABLE IF EXISTS fee_rates;
