# گزارش بازسازی پروژه — چه چیزی تغییر کرد

شاخه: `refactor/clean-arch-postgres-variants` (هر فاز یک commit معنادار).

## فاز ۰ — خط پایه
- حذف فراخوانی‌های `Ping` از asynq که با نسخهٔ فعلی کامپایل نمی‌شدند.

## فاز ۱ — بازساختار ساختاری به سبک gapbox (`refactor: restructure project to clean architecture layout`)
- انتقال کد از `internal/modules/{admin,public}` به چهار لایهٔ `domain / application / infrastructure / interfaces` + `bootstrap` + `pkg` + `cmd`.
- استخراج **اینترفیس** تمام ریپازیتوری‌ها به `domain/repositories` (رفع تداخل نام `AuthenticateRepositoryInterface` با پیشوند Admin/Customer).
- DTO ها به `application/dto/{admin,web}`، سرویس‌ها به `application/usecases/*`، پیاده‌سازی‌ها به `infrastructure/repositories/*`.
- `cmd` سبک gapbox: دستورات cobra (`serve`, `worker`, `scheduler`, `migrate`, `seed`) + `main.go` ریشه؛ حذف `cmd/http` و `cmd/commands`.
- HTML ها (بدون تغییر محتوا) به `templates/{admin,site,layouts,errors}` و استاتیک‌ها به `public/` منتقل شدند؛ `LoadHTMLGlob` و مسیرهای Static به‌روز شد.
- رفع باگ قدیمی: `InsertCart` مقدار تعداد آیتم تکراری را persist نمی‌کرد (`Save(&cart)` روی children اعمال نمی‌شود).

## فاز ۲ — PostgreSQL به‌جای MySQL (`feat: migrate database from MySQL to PostgreSQL 18 with goose migrations`)
- درایور `gorm.io/driver/postgres` + پکیج `infrastructure/database/postgres` (DSN با password نقل‌قولی تا مقدار خالی parsing را نشکند).
- **goose** جایگزین `sql-migrate` و `AutoMigrate` شد: فرمان‌های `migrate`, `migrate:status`, `migrate:rollback`, `migrate:reset`, `make:migration`.
- migration کامل ۲۱ جدول + FK + ایندکس‌های soft-delete (`20260923000000_create_shop_tables.sql`).
- `docker-compose` ها: `postgres:18-alpine` به‌جای mysql؛ افزودن سرویس Typesense به dev؛ حذف متغیرهای `MYSQL_*` و افزودن `POSTGRES_*`.
- حذف `gorm.io/driver/mysql` و `rubenv/sql-migrate` از go.mod.
- تست شده: `migrate up/down/reset` + seed کامل روی PostgreSQL محلی.

## فاز ۳ — حذف MongoDB → JSONB (`feat: remove MongoDB - store flattened product as JSONB read_model on products`)
- `products.read_model JSONB` + جدول `product_recommendations`.
- `SyncMongo` → **`SyncReadModel`**: همان مدل فلت (محصول + دسته + برند + تصاویر + ویژگی‌ها + موجودی‌ها با اتریبیوت) در JSONB ذخیره می‌شود؛ upsert تایپ‌سنس روی همان مسیر.
- صفحهٔ محصول فروشگاه از JSONB خوانده می‌شود (کلیدهای `_id` / `product` / `inventories` دقیقاً همان قرارداد قبلی)؛ `_id` حالا شناسهٔ **عددی** محصول است (نه ObjectID).
- سبد خرید: `product_id` عددی parse می‌شود؛ `InsertCart` دیگر `MongoProduct` نمی‌گیرد.
- پیشنهادها: `InsertRecommendation`/`GetAllRecommendation` روی جدول رابطه‌ای با اعتبارسنجی FK؛ پیشنهادهای صفحهٔ محصول از read_model ساخته می‌شوند.
- حذف کامل: پکیج mongodb، ریپازیتوری home_mongo، انتیتی/DTO های bson، فیلد MongoClient در bootstrap/events، سرویس‌های mongo در compose و env، درایور mongo از go.mod.
- تایپ‌سنس: گارد nil-client + health-check غیر-کشنده در بوت.

## فاز ۴ — مدل واریانت به سبک Laravel (`feat: Laravel-style variant model`)
- جداول جدید `product_variants` + `variant_attribute_values` با **حفظ شناسه‌ها** (migration داده + بازگشتی Down)؛ حذف `product_inventories` / `product_inventory_attributes`.
- `product_variants`: `stock`, `reserved_stock`, `price/sale_price/discount_price` (NULL = ارث از محصول)، `status` (`VariantStatusActive/Inactive`)، `expires_at`.
- ستون‌های aggregate + `attributes_json` (GIN) + `product_type` روی products؛ **backfill** در migration.
- **`application/usecases/pricing.PricingService`** → `RefreshProductAggregates`؛ بعد از هر عملیات ادمین روی واریانت و تغییر قیمت محصول صدا زده می‌شود؛ جریان سفارش (reserve/confirm) بعد از commit تراکنش refresh می‌کند (-sync خواندنی هم به بعد از commit منتقل شد — رفع باگ قدیمی همگام‌سازی زودهنگام).
- بازنویسی inventory repo روی واریانت‌ها (نگاشت استخر `product_attributes` → `attribute_value_id`، dedupe برای UNIQUE)؛ DTO های سفارش از `attribute_values` عنوان می‌خوانند؛ preload ها به `VariantAttributeValues.AttributeValue`.
- قرارداد HTML حفظ شد: کلیدهای `inventory_id/quantity/attributes`، alias ستون `stock AS quantity`، URL های `/admins/inventories/...` و `/admins/product-inventory-attributes/:id/delete`.
- سیدر: variant های تو در تو با `VariantAttributeValue` (بدون شناسه‌های hardcode)؛ کدهای `size`/`color`.
- تست شده: چرخه Down/Up goose، seed کامل (۲۰۶ واریانت)، محصول variable با `attributes_json={"size":[2,3,4]}`، رoundtrip PricingService، رندر صفحهٔ محصول با ۳ گزینه موجودی.

## فاز ۵ — اصلاح کامل Typesense
- **schema غنی**: `description`, `category`(facet), `brand`(facet), `original_price`, `sale_price`, `discount`, `stock`, `in_stock`(facet), `status`(facet) علاوه بر `id/title/slug/sku`.
- **sync**: آپلود کامل داکینگ هنگام `status=true`؛ **delete** از ایندکس هنگام انتشار-نکردن محصول (`status=false`) در همان `SyncReadModel`.
- **`search:reindex [--recreate]`**: بازسازی schema (drop+create) + نمایه‌سازی همه محصولات از read model ها (تست شد: `204 indexed, 0 failed`).
- هندلر جستجو: `query_by=title,sku,description,category,brand` با fallback به `title,sku` برای schema قدیمی؛ گارد nil client و hits؛ شکل پاسخ `{results:[{title,link}]}` برای `tsearch.html` ثابت است.

## فاز ۶ — تست‌نویسی (`test: unit, handler and DB-gated integration test suite`)
- unit: `domain_err.HandleError`، DTO های محصول/سفارش، `PricingService`، ابزارهای `pkg/util`.
- handler: پارس `product_id` عددی و رد کردن ObjectID قدیمی از طریق gin engine واقعی.
- integration (نیازمند `TEST_DATABASE_URL`، بدون نیاز به cleanup — rollback تراکنش): ماشین حساب aggregate ها (قیمت مؤثر/min-max/رزرو/variable/attributes_json)، قرارداد JSONB صفحه محصول، هوک‌های `AfterCreate`(attributes.code) و `BeforeCreate`(variant_attribute_values.product_id).
- `make test` / `make test-cover` — بدون DB تمیز skip می‌شود.

## فاز ۷ — داکیومنت
- `AGENTS.md` (قرارداد اجرا برای ایجنت‌ها)، `doc/architecture.md`، `doc/database.md`، این فایل، و به‌روزرسانی `README.md`.

## باگ‌های قدیمی که ضمن مسیر رفع شدند
1. `InsertCart` تعداد آیتم تکراری را persist نمی‌کرد.
2. `SyncReadModel` در `OrderPaidSuccessfully` قبل از commit تراکنش اجرا می‌شد (دادهٔ قدیمی می‌خواند).
3. health-check تایپ‌سنس در بوت `log.Fatal` بود (پایین بودن موتور جستجو کل سرور را می‌کشت).
4. `UpsertInTypesence` روی client صفر پانیک می‌کرد.
5. دستورات `start-worker`/`start-schedule` به فایل‌های main وجود-نداشته اشاره می‌کردند.

## محدودیت‌های عمدی
- **HTML تغییر نکرد** (فایل‌های `templates/` دست‌نخورده؛ جابجایی مسیر فقط در `LoadHTMLGlob`).
- قیمت per-variant از طریق فرم ادمین قابل ویرایش نیست (فرم HTML فقط quantity دارد)؛ ستون‌ها و محاسبات آماده‌اند و از کد/سیدر/API قابل تنظیم‌اند.
- نمایه‌سازی زندهٔ Typesense در این محیط تست نشد (موتور در دسترس نبود) — رفتار resilient (حذف/آپلود با گارد) و دستور reindex تست شدند.
- نمای دوتایی entity/model کامل gapbox (mappers) پیاده نشد؛ entity ها همچنان tag های GORM دارند (انحراف مستند از gapbox خالص).
