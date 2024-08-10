change add new products
<br/>
we need 3 page for this approach:
<br/>
1- add a product (images , videos, prices , title , slug ,sku , status )
<br/>
2- add a page for add attributes (product_attribute_values)
<br/>
3- add a page for add inventory for chosen attribute values

<br/>

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

<ul>
<li>add a page for product inventory  </li>
<li>use transaction for add attr-values for a product</li>
<li>implement upload media for a product</li>


</ul>

- add brand [title , slug ,image(just one image)]
- add brand_id to products table (nullable) [add its relation in product entity]
- check product images
- ability to upload videos of product (use product videos to upload
  videos [must use tusd pkg , fmtp to reduce video size])
- remove whitespace when u wanna insert any record to db [insert,update]
- implement show and edit product attributes
- implement show and edit product inventory
- impl edit inventory and edit product_attribute
- show attribute-value in add-attributes page

# todos
## 1. مدیریت محصولات
- [x] **افزودن/ویرایش/حذف محصولات:** امکان اضافه کردن محصولات جدید، ویرایش مشخصات محصولات (مانند قیمت، توضیحات، تصاویر)، و حذف محصولات قدیمی.
- [x] **مدیریت دسته‌بندی‌ها:** امکان ایجاد، ویرایش، و حذف دسته‌بندی‌های محصولات.
- **مدیریت موجودی:** مشاهده و مدیریت موجودی محصولات و ارسال اعلان‌های خودکار در صورت کمبود موجودی.
- **تنظیم قیمت و تخفیف:** امکان تنظیم قیمت، اعمال تخفیف‌های موقت یا دائمی، و نمایش قیمت قبلی و جدید.

## 2. مدیریت سفارشات
- **مشاهده و پیگیری سفارشات:** امکان مشاهده لیست کامل سفارشات با جزئیاتی مانند وضعیت پرداخت، وضعیت ارسال، و تاریخ ثبت سفارش.
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
