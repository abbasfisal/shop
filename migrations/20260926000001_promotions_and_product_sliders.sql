-- +goose Up
-- ============================================================
-- 1) promotion banners: title / schedule / layout (2-up | 4-up)
--    + ordering on the existing banners table
-- 2) homepage product sliders (max 10 products, position slot)
-- ============================================================

ALTER TABLE banners ADD COLUMN title       TEXT;
ALTER TABLE banners ADD COLUMN layout      VARCHAR(8) NOT NULL DEFAULT 'two';
ALTER TABLE banners ADD COLUMN starts_at   DATE;
ALTER TABLE banners ADD COLUMN ends_at     DATE;
ALTER TABLE banners ADD COLUMN sort_order  INTEGER NOT NULL DEFAULT 0;
ALTER TABLE banners ALTER COLUMN status SET DEFAULT TRUE;
UPDATE banners SET status = TRUE WHERE status IS NULL;
ALTER TABLE banners ALTER COLUMN status SET NOT NULL;

-- layout CHECK: only the two promotion shapes the design supports
ALTER TABLE banners ADD CONSTRAINT chk_banners_layout CHECK (layout IN ('two', 'four'));
CREATE INDEX idx_banners_promo ON banners (layout, status, sort_order);

-- ------------------------------------------------------------
-- product sliders
-- ------------------------------------------------------------
CREATE TABLE product_sliders (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    position    VARCHAR(32) NOT NULL DEFAULT 'after_categories',
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    status      VARCHAR(16) NOT NULL DEFAULT 'published'
                CHECK (status IN ('draft', 'published', 'archived')),
    starts_at   DATE,
    ends_at     DATE,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_product_sliders_deleted_at ON product_sliders (deleted_at);
CREATE INDEX idx_product_sliders_feed ON product_sliders (status, position, deleted_at);

COMMENT ON COLUMN product_sliders.position IS
    'homepage slot: after_slider | mid_content | after_categories | before_footer';
COMMENT ON COLUMN product_sliders.category_id IS
    'optional catalog scope for the "مشاهده همه" page';

CREATE TABLE slider_products (
    id         BIGSERIAL PRIMARY KEY,
    slider_id  BIGINT NOT NULL REFERENCES product_sliders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (slider_id, product_id)
);
CREATE INDEX idx_slider_products_deleted_at ON slider_products (deleted_at);
CREATE INDEX idx_slider_products_slider ON slider_products (slider_id, deleted_at);
CREATE INDEX idx_slider_products_product ON slider_products (product_id);

-- +goose Down
DROP TABLE IF EXISTS slider_products;
DROP TABLE IF EXISTS product_sliders;

DROP INDEX IF EXISTS idx_banners_promo;
ALTER TABLE banners DROP CONSTRAINT IF EXISTS chk_banners_layout;
ALTER TABLE banners DROP COLUMN IF EXISTS sort_order;
ALTER TABLE banners DROP COLUMN IF EXISTS ends_at;
ALTER TABLE banners DROP COLUMN IF EXISTS starts_at;
ALTER TABLE banners DROP COLUMN IF EXISTS layout;
ALTER TABLE banners DROP COLUMN IF EXISTS title;
ALTER TABLE banners ALTER COLUMN status DROP NOT NULL;
ALTER TABLE banners ALTER COLUMN status DROP DEFAULT;
