# دیتابیس — PostgreSQL (goose)

اتصال: متغیرهای `POSTGRES_*` در `.env` (نمونه: `.env.example`).
فرمان‌ها: `go run . migrate | migrate:status | migrate:rollback | migrate:reset | make:migration NAME=...`

## فایل‌های مایگریشن

| فایل | محتوا |
|---|---|
| `20260923000000_create_shop_tables.sql` | جدول‌های پایه: users, customers, addresses, brands, categories, attributes, attribute_values, products, product_images, features, product_attributes, product_inventories, product_inventory_attributes, carts, cart_items, orders, order_items, payments, otps, sessions, banners (+FK و ایندکس soft-delete) |
| `20260923030000_product_read_model.sql` | `products.read_model` (JSONB) + جدول `product_recommendations` (جایگزین کالکشن recommendations در MongoDB) |
| `20260924000000_product_variants.sql` | **مدل واریانت**: `product_variants` + `variant_attribute_values` + ستون‌های aggregate روی products (`product_type`, `min_price`, `max_price`, `total_stock`, `total_reserved`, `available_stock`, `in_stock`, `variants_count`, `attributes_json` + ایندکس GIN `idx_products_attrs_gin`) + `attributes.code/sort_order` + `attribute_values.sort_order`؛迁移 داده از جداول inventory قدیمی با **حفظ شناسه‌ها** (کلیدهای `cart_items.inventory_id` / `order_items.inventory_id` معتبر می‌مانند) و حذف جداول قدیمی (Down آن‌ها را بازسازی می‌کند) |

## نقشه جداول (بعد از migration سوم)

```
users, customers, addresses, banners, sessions, otps

categories (self-FK parent_id) ──┐
brands                           ├── products
                                 │     ├─ product_images
                                 │     ├─ features
                                 │     ├─ product_attributes (استخر اتریبیوت‌های تخصیص‌یافته — فرم add-attributes)
                                 │     ├─ product_variants (موجودی + قیمت هر ترکیب)
                                 │     │     └─ variant_attribute_values ── attribute_values ── attributes(code)
                                 │     ├─ read_model JSONB (مدل فلت storefront — جایگزین MongoDB)
                                 │     ├─ attributes_json JSONB + GIN (کش فیلتر: {"size":[2,3,4]})
                                 │     └─ aggregate: min/max_price, total_stock, total_reserved,
                                 │        available_stock, in_stock, variants_count, product_type
                                 └─ product_recommendations (پیشنهادها)

carts ── cart_items.inventory_id ──→ product_variants.id   (نام ستون حفظ شده، مقدار = variant)
orders ── order_items.inventory_id ──→ product_variants.id
```

## قواعد قیمت و موجودی (الگوی Laravel)

1. **قیمت مؤثر variant** = اگر `discount_price` معتبر بود (`>0` و `< قیمت مؤثر فروش`) از آن، وگرنه `variant.sale_price` اگر مقدار دارد، وگرنه `products.sale_price` (ارث).
2. **min_price / max_price** محصول = دامنه قیمت مؤثر واریانت‌های `active`.
3. **available_stock** = `max(total_stock - total_reserved, 0)`؛ `in_stock` = available > 0.
4. **product_type** = `variable` اگر محصول حداقل یک `variant_attribute_values` داشته باشد، وگرنه `simple`.
5. **attributes_json** از اتریبیوت‌های واریانت‌های active با کلید `attributes.code` (fallback `attr_<id>`) ساخته می‌شود.
6. بازمحاسبه بعد از **هر تغییر واریانت** توسط `application/usecases/pricing.PricingService.RefreshProductAggregates` (ادمین) و `product.RefreshProductAggregates` (جریان سفارش/seed در infrastructure — بعد از commit تراکنش).

## نکات مهاجرت داده

- شناسه‌های `product_inventories.id` عیناً به `product_variants.id` کپی شدند تا سوابق سبد/سفارش معتبر بمانند.
- لینک‌های قدیمی `product_inventory_attributes.id` هم حفظ شدند (ستون id در `variant_attribute_values`) تا URL های `delete_attribute_link` کار کنند.
- `attributes.code` برای ردیف‌های قدیمی با `attr_<id>` پر شد؛ سیدر از کد‌های `size`/`color` (همان نمونه Laravel) استفاده می‌کند.
