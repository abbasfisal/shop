-- +goose Up
-- ============================================================
-- Laravel-style variant model:
--   product_variants (per-combination stock + price)
--   variant_attribute_values (variant <-> attribute_value)
--   product aggregate cache columns + attributes_json (GIN)
--   attributes.code / sort_order
-- Old tables are migrated then dropped. Variant IDs are preserved
-- so cart_items.inventory_id / order_items.inventory_id stay valid.
-- ============================================================

CREATE TABLE product_variants (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id),
    sku             TEXT,
    price           BIGINT,          -- NULL = inherit products.original_price
    sale_price      BIGINT,          -- NULL = inherit products.sale_price
    discount_price  BIGINT,          -- golden rule: counts only when < effective sale price
    stock           BIGINT NOT NULL DEFAULT 0,
    reserved_stock  BIGINT NOT NULL DEFAULT 0,
    status          VARCHAR(16) NOT NULL DEFAULT 'active',
    expires_at      DATE,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);
CREATE INDEX idx_product_variants_deleted_at ON product_variants (deleted_at);

CREATE TABLE variant_attribute_values (
    id                 BIGSERIAL PRIMARY KEY,
    product_id         BIGINT NOT NULL REFERENCES products(id),
    variant_id         BIGINT NOT NULL REFERENCES product_variants(id),
    attribute_value_id BIGINT NOT NULL REFERENCES attribute_values(id),
    created_at         TIMESTAMPTZ DEFAULT now(),
    updated_at         TIMESTAMPTZ DEFAULT now(),
    deleted_at         TIMESTAMPTZ,
    UNIQUE (variant_id, attribute_value_id)
);
CREATE INDEX idx_vav_product_id ON variant_attribute_values (product_id);
CREATE INDEX idx_vav_variant_id ON variant_attribute_values (variant_id);

-- ---- migrate data: inventories -> variants (IDs preserved) ----
INSERT INTO product_variants (id, product_id, stock, reserved_stock, status, created_at, updated_at, deleted_at)
SELECT pi.id, pi.product_id, pi.quantity, pi.reserved_stock, 'active', pi.created_at, pi.updated_at, pi.deleted_at
FROM product_inventories pi;

SELECT setval(pg_get_serial_sequence('product_variants', 'id'),
              COALESCE((SELECT MAX(id) FROM product_variants), 0) + 1, false);

-- inventory-attribute links -> variant_attribute_values (IDs preserved,
-- dedup because the old table had no unique constraint)
INSERT INTO variant_attribute_values (id, product_id, variant_id, attribute_value_id, created_at, updated_at, deleted_at)
SELECT pia.id, pia.product_id, pia.product_inventory_id, pa.attribute_value_id,
       pia.created_at, pia.updated_at, pia.deleted_at
FROM product_inventory_attributes pia
JOIN product_attributes pa ON pa.id = pia.product_attribute_id
WHERE pia.deleted_at IS NULL
ON CONFLICT (variant_id, attribute_value_id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('variant_attribute_values', 'id'),
              COALESCE((SELECT MAX(id) FROM variant_attribute_values), 0) + 1, false);

DROP TABLE product_inventory_attributes;
DROP TABLE product_inventories;

-- ---- attributes: code + sort_order (Laravel pattern) ----
ALTER TABLE attributes ADD COLUMN code VARCHAR(64);
ALTER TABLE attributes ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;
UPDATE attributes SET code = 'attr_' || id::text WHERE code IS NULL OR code = '';
ALTER TABLE attributes ALTER COLUMN code SET NOT NULL;
CREATE UNIQUE INDEX idx_attributes_code ON attributes (code);

ALTER TABLE attribute_values ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

-- ---- products: type + aggregate cache + attributes_json ----
ALTER TABLE products ADD COLUMN product_type VARCHAR(16) NOT NULL DEFAULT 'simple';
ALTER TABLE products ADD COLUMN min_price BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN max_price BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN total_stock BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN total_reserved BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN available_stock BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN in_stock BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE products ADD COLUMN variants_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN attributes_json JSONB NOT NULL DEFAULT '{}';
CREATE INDEX idx_products_attrs_gin ON products USING GIN (attributes_json);

-- ---- backfill aggregates for existing rows ----
UPDATE products p SET
    total_stock = COALESCE(s.ts, 0),
    total_reserved = COALESCE(s.tr, 0),
    available_stock = GREATEST(COALESCE(s.ts, 0) - COALESCE(s.tr, 0), 0),
    in_stock = (GREATEST(COALESCE(s.ts, 0) - COALESCE(s.tr, 0), 0)) > 0,
    variants_count = COALESCE(s.vc, 0)
FROM (SELECT product_id,
             SUM(stock)              AS ts,
             SUM(reserved_stock)     AS tr,
             COUNT(*)                AS vc
      FROM product_variants
      WHERE deleted_at IS NULL AND status = 'active'
      GROUP BY product_id) s
WHERE p.id = s.product_id;

UPDATE products p SET
    min_price = COALESCE(e.minp, 0),
    max_price = COALESCE(e.maxp, 0)
FROM (SELECT v.product_id,
             MIN(CASE
                     WHEN v.discount_price IS NOT NULL AND v.discount_price > 0
                         AND v.discount_price < COALESCE(NULLIF(v.sale_price, 0), pr.sale_price)
                         THEN v.discount_price
                     ELSE COALESCE(NULLIF(v.sale_price, 0), pr.sale_price) END) AS minp,
             MAX(CASE
                     WHEN v.discount_price IS NOT NULL AND v.discount_price > 0
                         AND v.discount_price < COALESCE(NULLIF(v.sale_price, 0), pr.sale_price)
                         THEN v.discount_price
                     ELSE COALESCE(NULLIF(v.sale_price, 0), pr.sale_price) END) AS maxp
      FROM product_variants v
               JOIN products pr ON pr.id = v.product_id
      WHERE v.deleted_at IS NULL AND v.status = 'active'
      GROUP BY v.product_id) e
WHERE p.id = e.product_id;

UPDATE products p SET attributes_json = COALESCE((
    SELECT jsonb_object_agg(code, vids)
    FROM (SELECT a.code, jsonb_agg(DISTINCT vav.attribute_value_id) AS vids
          FROM variant_attribute_values vav
                   JOIN product_variants v ON v.id = vav.variant_id
                   JOIN attribute_values av ON av.id = vav.attribute_value_id
                   JOIN attributes a ON a.id = av.attribute_id
          WHERE v.product_id = p.id
            AND v.deleted_at IS NULL AND v.status = 'active'
            AND vav.deleted_at IS NULL AND av.deleted_at IS NULL AND a.deleted_at IS NULL
          GROUP BY a.code) sub), '{}'::jsonb);

UPDATE products p SET product_type = CASE
    WHEN EXISTS (SELECT 1
                 FROM variant_attribute_values vav
                          JOIN product_variants v ON v.id = vav.variant_id
                 WHERE v.product_id = p.id
                   AND v.deleted_at IS NULL AND v.status = 'active'
                   AND vav.deleted_at IS NULL)
        THEN 'variable'
    ELSE 'simple' END;

-- +goose Down
-- recreate legacy tables
CREATE TABLE product_inventories (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT NOT NULL REFERENCES products(id),
    quantity        BIGINT NOT NULL DEFAULT 0,
    reserved_stock  BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX idx_product_inventories_product_id ON product_inventories (product_id);
CREATE INDEX idx_product_inventories_deleted_at ON product_inventories (deleted_at);

CREATE TABLE product_inventory_attributes (
    id                    BIGSERIAL PRIMARY KEY,
    product_id            BIGINT NOT NULL REFERENCES products(id),
    product_inventory_id  BIGINT NOT NULL REFERENCES product_inventories(id),
    product_attribute_id  BIGINT NOT NULL REFERENCES product_attributes(id),
    created_at            TIMESTAMPTZ DEFAULT now(),
    updated_at            TIMESTAMPTZ DEFAULT now(),
    deleted_at            TIMESTAMPTZ
);
CREATE INDEX idx_pia_product_id ON product_inventory_attributes (product_id);
CREATE INDEX idx_pia_product_inventory_id ON product_inventory_attributes (product_inventory_id);
CREATE INDEX idx_pia_product_attribute_id ON product_inventory_attributes (product_attribute_id);
CREATE INDEX idx_pia_deleted_at ON product_inventory_attributes (deleted_at);

INSERT INTO product_inventories (id, product_id, quantity, reserved_stock, created_at, updated_at, deleted_at)
SELECT id, product_id, stock, reserved_stock, created_at, updated_at, deleted_at
FROM product_variants;
SELECT setval(pg_get_serial_sequence('product_inventories', 'id'),
              COALESCE((SELECT MAX(id) FROM product_inventories), 0) + 1, false);

INSERT INTO product_inventory_attributes (id, product_id, product_inventory_id, product_attribute_id, created_at, updated_at, deleted_at)
SELECT vav.id, vav.product_id, vav.variant_id, pa.id, vav.created_at, vav.updated_at, vav.deleted_at
FROM variant_attribute_values vav
         JOIN product_attributes pa
              ON pa.product_id = vav.product_id AND pa.attribute_value_id = vav.attribute_value_id
WHERE vav.deleted_at IS NULL
ON CONFLICT DO NOTHING;
SELECT setval(pg_get_serial_sequence('product_inventory_attributes', 'id'),
              COALESCE((SELECT MAX(id) FROM product_inventory_attributes), 0) + 1, false);

DROP TABLE variant_attribute_values;
DROP TABLE product_variants;

DROP INDEX idx_products_attrs_gin;
ALTER TABLE products DROP COLUMN product_type;
ALTER TABLE products DROP COLUMN min_price;
ALTER TABLE products DROP COLUMN max_price;
ALTER TABLE products DROP COLUMN total_stock;
ALTER TABLE products DROP COLUMN total_reserved;
ALTER TABLE products DROP COLUMN available_stock;
ALTER TABLE products DROP COLUMN in_stock;
ALTER TABLE products DROP COLUMN variants_count;
ALTER TABLE products DROP COLUMN attributes_json;

DROP INDEX idx_attributes_code;
ALTER TABLE attributes DROP COLUMN code;
ALTER TABLE attributes DROP COLUMN sort_order;
ALTER TABLE attribute_values DROP COLUMN sort_order;
