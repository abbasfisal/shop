# معماری پروژه (Clean Architecture — الگوی gapbox)

## لایه‌ها و جهت وابستگی

```
interfaces (HTTP/worker)
    ↓
application (usecases / dto / pricing)
    ↓
infrastructure (repositories / database / events)
    ↓
domain (entities / repository interfaces / domain_err)
```

| لایه | مسیر | مسئولیت |
|---|---|---|
| **domain** | `domain/entities` | انتیتی‌های خالص دامنه (Product, ProductVariant, Attribute, …) |
| | `domain/repositories` | **اینترفیس** ریپازیتوری‌ها (پیاده‌سازی در infrastructure) |
| | `domain/domain_err` | خطاهای sentinel با پیام فارسی |
| **application** | `application/usecases/*` | سرویس‌های کسب‌وکار (admin: product/category/… ، web: home) |
| | `application/dto/{admin,web}` | ساختار دادهٔ خروجی به تمپلیت‌ها/JSON |
| | `application/usecases/pricing` | **PricingService** — بازمحاسبهٔ aggregate های محصول |
| **infrastructure** | `infrastructure/repositories/*` | پیاده‌سازی اینترفیس‌ها با GORM/PostgreSQL |
| | `infrastructure/database/postgres` | اتصال PostgreSQL (singleton) |
| | `infrastructure/database/typesenceclient` | کلاینت Typesense |
| | `infrastructure/seeders` | دیتای نمونه |
| | `infrastructure/events` | Event/Listener (asynq) |
| **interfaces** | `interfaces/http/handlers/{admin,web}` | هندلرهای Gin |
| | `interfaces/http/routes` | تعریف مسیرها + ثبت middleware ها |
| | `interfaces/http/{requests,response,middleware}` | binding، رندر HTML، middleware |
| | `interfaces/worker/tasks` | تسک‌های asynq (admin jobs) |
| **bootstrap** | `bootstrap/` | بارگذاری config (.env + viper) و ساخت `Dependencies` |
| **cmd** | `cmd/` | دستورات cobra: serve / worker / scheduler / migrate / seed / search:reindex |
| **pkg** | `pkg/*` | ابزارهای عمومی (cache, sessions, helpers, util, …) |

## جریان یک درخواست

```
HTTP → Gin middleware (interfaces/http/routes)
    → Handler (interfaces/http/handlers)
    → Usecase (application/usecases)
    → Repository interface (domain/repositories)
    → Repository impl (infrastructure/repositories) → PostgreSQL
    → DTO (application/dto) → template (templates/) یا JSON
```

DI به‌صورت **دستی** در `interfaces/http/routes` + `bootstrap.Initialize()` انجام می‌شود (بدون فریمورک DI).

## تصمیم‌های کلیدی

1. **فقط PostgreSQL** — مایگریشن‌ها goose SQL هستند (`migrations/*.sql`)؛ AutoMigrate حذف شده.
2. **بدون MongoDB** — مدل فلت محصول در `products.read_model` (JSONB) نگهداری می‌شود و با `SyncReadModel` بعد از هر تغییر بازسازی می‌شود.
3. **مدل واریانت (الگوی Laravel)** — موجودی/قیمت هر ترکیب اتریبیوت در `product_variants`؛ لینک اتریبیوت‌ها در `variant_attribute_values`؛ خلاصه‌ٔ فیلتر در `products.attributes_json` (ایندکس GIN)؛ قیمت variant وقتی NULL باشد از محصول ارث می‌برد.
4. **HTML ثابت** — فایل‌های `templates/` تغییر نمی‌کنند؛ قرارداد آن‌ها در `AGENTS.md` مستند شده.
5. **Typesense اختیارتی در بوت** — اگر در دسترس نباشد فقط warning می‌زند؛ ایندکس غنی با `search:reindex` قابل بازسازی است.
