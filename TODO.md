explain WITH RECURSIVE CategoryHierarchy AS (
SELECT id, title, parent_id
FROM categories
WHERE id = (SELECT category_id FROM products WHERE id = 1)

    UNION ALL

    SELECT c.id, c.title, c.parent_id
    FROM categories c
             INNER JOIN CategoryHierarchy ch ON c.id = ch.parent_id

)
SELECT *
FROM CategoryHierarchy
WHERE parent_id IS NULL
LIMIT 1;


<br/>
- ability to upload videos of product (use product videos to upload
  videos [must use tusd pkg , fmtp to reduce video size])
- use record lock when editing or delete product records

# todos

## 1. مدیریت محصولات

- [x] **افزودن/ویرایش/حذف محصولات:** امکان اضافه کردن محصولات جدید، ویرایش مشخصات محصولات (مانند قیمت، توضیحات، تصاویر)،
  و حذف محصولات قدیمی.
- [x] **مدیریت دسته‌بندی‌ها:** امکان ایجاد، ویرایش، و حذف دسته‌بندی‌های محصولات.
- **مدیریت موجودی:** مشاهده و مدیریت موجودی محصولات و ارسال اعلان‌های خودکار در صورت کمبود موجودی.
- **تنظیم قیمت و تخفیف:** امکان تنظیم قیمت، اعمال تخفیف‌های موقت یا دائمی، و نمایش قیمت قبلی و جدید.

## 2. مدیریت سفارشات

- **مشاهده و پیگیری سفارشات:** امکان مشاهده لیست کامل سفارشات با جزئیاتی مانند وضعیت پرداخت، وضعیت ارسال، و تاریخ ثبت
  سفارش.
- **تایید/رد سفارشات:** امکان تایید یا رد سفارشات و ارسال اعلان به کاربر.
- **مدیریت مرجوعی‌ها:** سیستم مدیریت درخواست‌های مرجوعی و رسیدگی به آن‌ها.
- **تولید فاکتور:** امکان تولید و ارسال فاکتور به مشتریان.

## 3. مدیریت کاربران

- **مدیریت مشتریان:** مشاهده پروفایل‌ها و تاریخچه خرید مشتریان، امکان قفل یا حذف حساب‌های کاربری.
- **مدیریت دسترسی‌ها:** امکان ایجاد، ویرایش، و حذف نقش‌ها و سطوح دسترسی مختلف برای ادمین‌ها.
- **مدیریت نظرات و بازخوردها:** مشاهده و مدیریت نظرات کاربران در مورد محصولات، امکان تایید یا حذف نظرات.

## 4. گزارشات و تحلیل‌ها

- **گزارش فروش:** نمایش گزارشات جامع از فروش‌ها به صورت دوره‌ای (روزانه، هفتگی، ماهانه).
- **تحلیل رفتار مشتری:** نمایش تحلیل‌هایی از رفتار مشتریان مانند پرفروش‌ترین محصولات، محصولات پر بازدید، و نرخ تبدیل.
- **گزارش مالی:** نمایش گزارشات مالی شامل درآمد، هزینه‌ها، سود و زیان.

## 5. مدیریت محتوا و تبلیغات

- **مدیریت بنرها و تبلیغات:** امکان اضافه کردن و مدیریت بنرها و تبلیغات در سایت.
- **مدیریت صفحات و بلاگ:** ایجاد و ویرایش صفحات اطلاعاتی، مقالات و محتواهای وبلاگ.
- **مدیریت SEO:** تنظیمات سئو برای صفحات مختلف به منظور بهبود رتبه سایت در موتورهای جستجو.

## 6. پشتیبانی و ارتباطات

- **مدیریت تیکت‌ها و درخواست‌های پشتیبانی:** امکان مشاهده و پاسخگویی به درخواست‌های پشتیبانی کاربران.
- **ارسال ایمیل‌های انبوه:** امکان ارسال ایمیل‌های تبلیغاتی یا اطلاع‌رسانی به گروهی از کاربران.
- **تنظیم اعلان‌ها:** مدیریت اعلان‌ها و پیام‌های ارسالی به کاربران از طریق ایمیل، پیامک، یا سیستم اعلان داخلی.

  # 📝
- ارسال نوتیفیکیشن کاهش موجودی یا اتمام موجودی
- غیر فعال شدن و یا تگ اتمام موجودی برای محصول بدون موجودی

___

- [x] api for delete product inventory attribute
- [x] api for delete inventory and its attributes in the product_inventory_attributes table
- [x] Add attributes to an existing inventory
- [x] update inventory quantity
- [x] add home page html
- [x] add login / register html form
- [x] implement login / register for customer
- [x] send otp
- [x] resend otp
- [x] add customer profile html form
- [x] add customer address html form
- [ ] add product recom by elastic
- [x] add logOut
- [x] edit customer profile (name,last name)
- [x] customer profile manager
- [x] show all customer in admin panel
- [x] implement redis as a cache driver
- [ ] implement cache system for menu and rows(most orders,most viewed , newest , by categories , etc)
- [x] add enable/disable , priority to category
- [x] store final product in mongodb(flat)
- [x] crud product feature
- [x] show single product in home page with multi qty
- [x] show product by category in home page

__
- [x] separate SyncMongo() from product_repository
- [x] add some real menu
- [x] load all menu in redis
- 
- [x] add some real product
- [x] store product with all attributes in mongodb
- load data rows(new , most ordered , ) in home page
- ability to bookmark product by user
- log all request in elastic
- order 
- [x] get tax-code
- [x] buy product by bank gateway
- [x] show customer paginated orders
- [x] show detail of an order in customer side
- [x] show paginated orders in admin panel
- [x] show detail of an order in admin side 
- [x] implement event/listener system
- [x] implement job/queue/schedule system using asynq pkg
- [x] install [Asynqmon](https://github.com/hibiken/asynqmon) Web UI for monitoring & administering Asynq task queue

add make file :

`docker-compose -f docker-compose.dev.yml --env-file .env.development up
`



-----------_/\_------------
`sample rediLock code`
```go
package services

import (
	"context"
	"fmt"
	"time"

	"project/internal/models"
	"project/internal/repositories"
	"project/pkg/redis"

	"github.com/redis/go-redis/v9"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, order *models.Order) error
}

type orderService struct {
	orderRepo repositories.OrderRepository
	redis     *redis.Client
}

func NewOrderService(orderRepo repositories.OrderRepository, redisClient *redis.Client) OrderService {
	return &orderService{orderRepo: orderRepo, redis: redisClient}
}

func (s *orderService) PlaceOrder(ctx context.Context, order *models.Order) error {
	// قفل‌ها برای محصولات
	lockKeys := make([]string, 0)
	for _, item := range order.Items {
		lockKey := fmt.Sprintf("lock:product:%d", item.ProductID)
		locked, err := s.redis.SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
		if err != nil {
			// آزاد کردن قفل‌های قبلی در صورت خطا
			s.releaseLocks(ctx, lockKeys)
			return fmt.Errorf("error acquiring lock for product %d: %w", item.ProductID, err)
		}
		if !locked {
			s.releaseLocks(ctx, lockKeys)
			return fmt.Errorf("product %d is locked by another process", item.ProductID)
		}
		lockKeys = append(lockKeys, lockKey)
	}

	// آزاد کردن قفل‌ها در انتها
	defer s.releaseLocks(ctx, lockKeys)

	// ذخیره سفارش و آیتم‌ها در دیتابیس
	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	return nil
}

func (s *orderService) releaseLocks(ctx context.Context, lockKeys []string) {
	for _, lockKey := range lockKeys {
		s.redis.Del(ctx, lockKey)
	}
}


```