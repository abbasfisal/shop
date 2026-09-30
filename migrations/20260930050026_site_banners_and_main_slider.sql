-- +goose Up
-- ============================================================
-- 1) site_banners — بنرهای ثابت سایت (جایگاه: هدر / آینده: اسلایدر تصویری)
--    تصویر دسکتاپ اجباری + تصویر موبایل اختیاری (در نبود، از دسکتاپ استفاده می‌شود)
-- 2) product_sliders.subtitle — زیرعنوان اسلایدر (زیر تیتر در اسلایدر اصلی)
-- ============================================================

CREATE TABLE site_banners (
    id           BIGSERIAL PRIMARY KEY,
    placement    VARCHAR(16) NOT NULL DEFAULT 'header'
                 CHECK (placement IN ('header', 'main')),
    title        TEXT,
    link         TEXT,
    image        TEXT NOT NULL,
    mobile_image TEXT,
    status       BOOLEAN NOT NULL DEFAULT TRUE,
    starts_at    DATE,
    ends_at      DATE,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_site_banners_deleted_at ON site_banners (deleted_at);
-- feed index: فید فروشگاه (جایگاه + فعال + ترتیب)
CREATE INDEX idx_site_banners_feed ON site_banners (placement, status, sort_order, deleted_at);

COMMENT ON TABLE site_banners IS
    'بنرهای ثابت سایت (هدر و جایگاه‌های تصویری آتی) — جدا از بنرهای پروموشن جدول banners';
COMMENT ON COLUMN site_banners.placement IS
    'header = بنر بالای سایت | main = رزرو شده برای بنر تصویری اسلایدر اصلی';
COMMENT ON COLUMN site_banners.image IS
    'تصویر نسخه دسکتاپ (الزامی)';
COMMENT ON COLUMN site_banners.mobile_image IS
    'تصویر نسخه موبایل (اختیاری — در خالی بودن از image استفاده می‌شود)';

-- ------------------------------------------------------------
-- اسلایدر اصلی: زیرعنوان (مثل «بر اساس سلیقه شما»)
-- ------------------------------------------------------------
ALTER TABLE product_sliders ADD COLUMN subtitle TEXT;
COMMENT ON COLUMN product_sliders.subtitle IS
    'زیرعنوان زیر تیتر اسلایدر (اختیاری) — در اسلایدر اصلی صفحه نمایش داده می‌شود';

-- +goose Down
ALTER TABLE product_sliders DROP COLUMN IF EXISTS subtitle;
DROP TABLE IF EXISTS site_banners;
