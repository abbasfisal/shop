package seeders

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"shop/domain/entities"
	productRepo "shop/infrastructure/repositories/product"
)

// ============================================================
// Demo seeder
//
// Seeds the content introduced by the recent shop features:
//
//   - products in every supported state: simple / variable,
//     published / draft / archived, with / without discount,
//     in / out of stock, inactive variant, expiry dates
//   - promotion banners (two-up / four-up) with status + schedule
//   - homepage product sliders (one per position, max 10 products)
//
// Every part is idempotent: a re-run reports "already seeded" and
// inserts nothing, so `go run . seed` is safe to execute twice.
// ============================================================

const demoPrefix = "DEMO-"
const demoBannerTitlePrefix = "پروموشن نمونه"
const demoSliderSlugPrefix = "demo-"

// theme assets copied into the banner upload folder
const (
	bannerAssetDir  = "public/shop/img/banner"
	bannerUploadDir = "public/uploads/media/banners/seed"
	productImageDir = "2024/09/27"
)

// SeedDemoShop is called by Seed() right before the read models are built.
func SeedDemoShop(db *gorm.DB) {
	seedDemoProducts(db)
	seedPromoBanners(db)
	seedDemoSliders(db)
	seedFeeRates(db)
}

// ------------------------------------------------------------
// 1. products — every supported state
// ------------------------------------------------------------

type demoVariant struct {
	Price         *uint // nil → inherit products.original_price
	SalePrice     *uint // nil → inherit products.sale_price
	DiscountPrice *uint
	Stock         uint
	Status        string
	ExpiresAt     *time.Time
	Values        []string // "attributeCode:value" pairs
}

type demoProduct struct {
	SKU           string
	Title         string
	Slug          string
	Description   string
	CategorySlug  string
	BrandSlug     string
	Status        string
	OriginalPrice uint
	SalePrice     uint
	ExpiresAt     *time.Time
	Images        []string
	Variants      []demoVariant
}

func demoProducts(now time.Time) []demoProduct {
	ptr := func(t time.Time) *time.Time { return &t }
	past := ptr(now.AddDate(0, 0, -5))
	future := ptr(now.AddDate(0, 0, 45))

	uintPtr := func(v uint) *uint { return &v }

	return []demoProduct{
		// --- simple · published · discounted · in stock ---
		{
			SKU:           demoPrefix + "SIMPLE-01",
			Title:         "تیشرت نایک کلاسیک مردانه",
			Slug:          "demo-nike-classic-tshirt",
			Description:   "تیشرت نخی مردانه با الگوی راحت، مناسب استفاده روزمره و ورزش.",
			CategorySlug:  "category-men-tee-shirts",
			BrandSlug:     "zara",
			Status:        entities.ProductStatusPublished,
			OriginalPrice: 450_000,
			SalePrice:     390_000,
			Images:        []string{"1.webp", "2.webp"},
			Variants: []demoVariant{
				// NULL price columns → inherit the product price (golden rule)
				{DiscountPrice: uintPtr(350_000), Stock: 25, Status: entities.VariantStatusActive},
			},
		},

		// --- simple · published · no discount · out of stock ---
		{
			SKU:           demoPrefix + "SIMPLE-02",
			Title:         "شلوار جین مردانه لی‌اس",
			Slug:          "demo-mens-jeans",
			Description:   "شلوار جین کلاسیک با پارچه کشسان و دوخت تمیز.",
			CategorySlug:  "category-men-trousers-jumpsuits",
			BrandSlug:     "iran",
			Status:        entities.ProductStatusPublished,
			OriginalPrice: 900_000,
			SalePrice:     780_000,
			Images:        []string{"3.webp"},
			Variants: []demoVariant{
				{Stock: 0, Status: entities.VariantStatusActive}, // ناموجود
			},
		},

		// --- simple · draft (پیش‌نویس — در فروشگاه دیده نمی‌شود) ---
		{
			SKU:           demoPrefix + "SIMPLE-03",
			Title:         "کیف دوشی چرم طبیعی",
			Slug:          "demo-leather-bag",
			Description:   "کیف دوشی چرم طبیعی با بند قابل تنظیم.",
			CategorySlug:  "category-men-accessories",
			BrandSlug:     "zara",
			Status:        entities.ProductStatusDraft,
			OriginalPrice: 1_200_000,
			SalePrice:     1_050_000,
			Images:        []string{"11.webp"},
			Variants: []demoVariant{
				{Stock: 12, Status: entities.VariantStatusActive},
			},
		},

		// --- simple · archived · product expiry in the past ---
		{
			SKU:           demoPrefix + "SIMPLE-04",
			Title:         "کاپشن زمستانی مردانه",
			Slug:          "demo-winter-jacket",
			Description:   "کاپشن زمستانی ضدآب با آستر داخلی گرم.",
			CategorySlug:  "category-men-sweatshirts",
			BrandSlug:     "bailando",
			Status:        entities.ProductStatusArchived,
			OriginalPrice: 2_500_000,
			SalePrice:     1_980_000,
			ExpiresAt:     past,
			Images:        []string{"22.webp"},
			Variants: []demoVariant{
				{Stock: 6, Status: entities.VariantStatusActive, ExpiresAt: past},
			},
		},

		// --- variable · published · mixed discounts & stock · one inactive ---
		{
			SKU:           demoPrefix + "VAR-01",
			Title:         "کفش ورزشی آدیداس ران‌فالکن",
			Slug:          "demo-adidas-runfalcon",
			Description:   "کفش دویدن آدیداس با زیره انعطاف‌پذیر و رویه مشبک.",
			CategorySlug:  "category-men-tee-shirts",
			BrandSlug:     "bailando",
			Status:        entities.ProductStatusPublished,
			OriginalPrice: 3_200_000,
			SalePrice:     2_900_000,
			Images:        []string{"33.webp", "111.webp", "222.webp"},
			Variants: []demoVariant{
				{Price: uintPtr(3_200_000), SalePrice: uintPtr(2_900_000),
					DiscountPrice: uintPtr(2_610_000), Stock: 8, Status: entities.VariantStatusActive,
					Values: []string{"size:s", "color:آبی"}},
				{Price: uintPtr(3_200_000), SalePrice: uintPtr(2_900_000), Stock: 5,
					Status: entities.VariantStatusActive, Values: []string{"size:m", "color:آبی"}},
				{Price: uintPtr(3_400_000), SalePrice: uintPtr(3_050_000),
					DiscountPrice: uintPtr(2_800_000), Stock: 0, Status: entities.VariantStatusActive,
					Values: []string{"size:s", "color:قرمز"}}, // ناموجود
				{Price: uintPtr(3_400_000), SalePrice: uintPtr(3_050_000), Stock: 4,
					Status: entities.VariantStatusInactive, Values: []string{"size:l", "color:قرمز"}},
			},
		},

		// --- variable · published · future expiry on one combo ---
		{
			SKU:           demoPrefix + "VAR-02",
			Title:         "مانتو مجلسی زنانه",
			Slug:          "demo-womens-formal-manteau",
			Description:   "مانتو مجلسی زنانه با پارچه لطیف و دوخت ظریف.",
			CategorySlug:  "category-women-shirts",
			BrandSlug:     "woody-sence",
			Status:        entities.ProductStatusPublished,
			OriginalPrice: 1_750_000,
			SalePrice:     1_590_000,
			ExpiresAt:     future,
			Images:        []string{"1.webp", "3.webp"},
			Variants: []demoVariant{
				{Stock: 7, Status: entities.VariantStatusActive, Values: []string{"size:m", "color:بنفش"}},
				{DiscountPrice: uintPtr(1_450_000), Stock: 3, Status: entities.VariantStatusActive,
					ExpiresAt: future, Values: []string{"size:l", "color:بنفش"}},
				{Stock: 9, Status: entities.VariantStatusActive, Values: []string{"size:m", "color:قرمز"}},
			},
		},

		// --- variable · draft ---
		{
			SKU:           demoPrefix + "VAR-03",
			Title:         "هودی مردانه اورجینال",
			Slug:          "demo-mens-hoodie",
			Description:   "هودی مردانه با جنس داخلی کرکی.",
			CategorySlug:  "category-men-hoodies",
			BrandSlug:     "iran",
			Status:        entities.ProductStatusDraft,
			OriginalPrice: 1_100_000,
			SalePrice:     990_000,
			Images:        []string{"11.webp", "22.webp"},
			Variants: []demoVariant{
				{Stock: 10, Status: entities.VariantStatusActive, Values: []string{"size:s", "color:آبی"}},
				{Stock: 6, Status: entities.VariantStatusActive, Values: []string{"size:l", "color:آبی"}},
			},
		},
	}
}

func seedDemoProducts(db *gorm.DB) {
	var count int64
	db.Model(&entities.Product{}).Where("sku LIKE ?", demoPrefix+"%").Count(&count)
	if count > 0 {
		fmt.Printf("[seed] demo products ........ already seeded (%d)\n", count)
		return
	}

	now := time.Now()
	ctx := context.Background()

	for _, spec := range demoProducts(now) {
		product := entities.Product{
			Sku:           spec.SKU,
			Title:         spec.Title,
			Slug:          spec.Slug,
			Description:   spec.Description,
			CategoryID:    demoCategoryID(db, spec.CategorySlug),
			BrandID:       demoBrandID(db, spec.BrandSlug),
			Status:        spec.Status,
			OriginalPrice: spec.OriginalPrice,
			SalePrice:     spec.SalePrice,
			ExpiresAt:     spec.ExpiresAt,
		}
		for _, image := range spec.Images {
			product.ProductImages = append(product.ProductImages, &entities.ProductImages{Path: productImageDir + "/" + image})
		}

		if err := db.Create(&product).Error; err != nil {
			fmt.Printf("[seed] demo product %s failed: %v\n", spec.SKU, err)
			continue
		}

		for _, vs := range spec.Variants {
			variant := entities.ProductVariant{
				ProductID:     product.ID,
				Price:         vs.Price,
				SalePrice:     vs.SalePrice,
				DiscountPrice: vs.DiscountPrice,
				Stock:         vs.Stock,
				Status:        vs.Status,
				ExpiresAt:     vs.ExpiresAt,
			}
			if variant.Status == "" {
				variant.Status = entities.VariantStatusActive
			}
			if err := db.Create(&variant).Error; err != nil {
				fmt.Printf("[seed] demo variant for %s failed: %v\n", spec.SKU, err)
				continue
			}

			for _, pair := range vs.Values {
				valueID := demoAttrValueID(db, pair)
				if valueID == 0 {
					fmt.Printf("[seed]   unknown attribute value %q (skipped)\n", pair)
					continue
				}
				link := entities.VariantAttributeValue{
					ProductID:        product.ID,
					VariantID:        variant.ID,
					AttributeValueID: valueID,
				}
				if err := db.Create(&link).Error; err != nil {
					fmt.Printf("[seed]   variant link %q failed: %v\n", pair, err)
				}
			}
		}

		// aggregates + read model + Typesense for this product only
		if err := productRepo.RefreshProductAggregates(ctx, db, product.ID); err != nil {
			fmt.Printf("[seed]   pricing refresh failed for %s: %v\n", spec.SKU, err)
		}
		if err := productRepo.SyncReadModel(ctx, db, product.ID); err != nil {
			fmt.Printf("[seed]   read model sync failed for %s: %v\n", spec.SKU, err)
		}
	}

	fmt.Println("[seed] demo products ......... done (simple/variable × published/draft/archived)")
}

// ------------------------------------------------------------
// 2. promotion banners (two-up / four-up + status + schedule)
// ------------------------------------------------------------

type demoBanner struct {
	Title     string
	Layout    string
	Link      string
	Status    bool
	StartsAt  *time.Time
	EndsAt    *time.Time
	SortOrder int
	Image     string // file name inside bannerAssetDir
}

func demoBanners(now time.Time) []demoBanner {
	ptr := func(t time.Time) *time.Time { return &t }
	started := ptr(now.AddDate(0, 0, -10))
	endsSoon := ptr(now.AddDate(0, 0, 60))
	expired := ptr(now.AddDate(0, 0, -3))

	return []demoBanner{
		// two-up row
		{Title: demoBannerTitlePrefix + " — حراج پاییز", Layout: entities.BannerLayoutTwo,
			Link: "/search/apparel", Status: true, SortOrder: 1, Image: "medium-banner-1.jpg"},
		{Title: demoBannerTitlePrefix + " — جشنواره بهار", Layout: entities.BannerLayoutTwo,
			Link: "/", Status: true, StartsAt: started, EndsAt: endsSoon,
			SortOrder: 2, Image: "medium-banner-2.jpg"},

		// four-up row
		{Title: demoBannerTitlePrefix + " — تخفیف شماره ۱", Layout: entities.BannerLayoutFour,
			Link: "/search/category-men-clothing", Status: true, SortOrder: 3, Image: "small-banner-1.jpg"},
		{Title: demoBannerTitlePrefix + " — تخفیف شماره ۲", Layout: entities.BannerLayoutFour,
			Link: "/search/category-women-shirts", Status: true, SortOrder: 4, Image: "small-banner-2.jpg"},

		// غیرفعال → در فروشگاه نمایش داده نمی‌شود
		{Title: demoBannerTitlePrefix + " — تخفیف شماره ۳ (غیرفعال)", Layout: entities.BannerLayoutFour,
			Status: false, SortOrder: 5, Image: "small-banner-3.jpg"},

		// منقضی → بازه زمانی تمام شده
		{Title: demoBannerTitlePrefix + " — تخفیف شماره ۴ (منقضی)", Layout: entities.BannerLayoutFour,
			Status: true, EndsAt: expired, SortOrder: 6, Image: "small-banner-4.jpg"},
	}
}

func seedPromoBanners(db *gorm.DB) {
	var count int64
	db.Model(&entities.Banner{}).Where("title LIKE ?", demoBannerTitlePrefix+"%").Count(&count)
	if count > 0 {
		fmt.Printf("[seed] promo banners ......... already seeded (%d)\n", count)
		return
	}

	now := time.Now()
	for _, spec := range demoBanners(now) {
		imageName, err := copyBannerAsset(spec.Image)
		if err != nil {
			fmt.Printf("[seed] banner image %s failed: %v\n", spec.Image, err)
			continue
		}

		banner := entities.Banner{
			Title:     spec.Title,
			Layout:    spec.Layout,
			Link:      spec.Link,
			Status:    spec.Status,
			StartsAt:  spec.StartsAt,
			EndsAt:    spec.EndsAt,
			SortOrder: spec.SortOrder,
			Priority:  uint(spec.SortOrder),
			Image:     imageName,
		}
		if banner.Layout == "" {
			banner.Layout = entities.BannerLayoutTwo
		}
		// Select forces every listed column: without it GORM drops zero
		// values that carry a default tag (status=false would become true).
		if err := db.
			Select("Title", "Layout", "Link", "Status", "StartsAt", "EndsAt",
				"SortOrder", "Priority", "Image").
			Create(&banner).Error; err != nil {
			fmt.Printf("[seed] banner %q failed: %v\n", spec.Title, err)
		}
	}

	fmt.Println("[seed] promo banners ......... done (two-up/four-up × active/inactive/scheduled/expired)")
}

// ------------------------------------------------------------
// 3. homepage product sliders (one per position)
// ------------------------------------------------------------

type demoSlider struct {
	Title        string
	Slug         string
	Position     string
	CategorySlug string
	Status       string
	StartsAt     *time.Time
	EndsAt       *time.Time
	Products     int
	PreferDemo   bool
}

func demoSliders(now time.Time) []demoSlider {
	ptr := func(t time.Time) *time.Time { return &t }
	started := ptr(now.AddDate(0, 0, -3))
	endsSoon := ptr(now.AddDate(0, 0, 90))

	return []demoSlider{
		{Title: "جدیدترین محصولات", Slug: demoSliderSlugPrefix + "newest",
			Position: entities.SliderPositionAfterSlider,
			Status:   entities.ProductStatusPublished, Products: 6, PreferDemo: true},

		{Title: "پیشنهاد ویژه", Slug: demoSliderSlugPrefix + "special-offer",
			Position: entities.SliderPositionMidContent,
			Status:   entities.ProductStatusPublished, Products: 5, PreferDemo: true},

		{Title: "منتخب پوشاک مردانه", Slug: demoSliderSlugPrefix + "mens-picks",
			Position:     entities.SliderPositionAfterCategories,
			CategorySlug: "category-men-clothing",
			Status:       entities.ProductStatusPublished, Products: 8, PreferDemo: true},

		{Title: "حراج پایان فصل", Slug: demoSliderSlugPrefix + "clearance",
			Position: entities.SliderPositionBeforeFooter,
			Status:   entities.ProductStatusPublished, StartsAt: started, EndsAt: endsSoon,
			Products: 5, PreferDemo: false},

		// پیش‌نویس → در فروشگاه نمایش داده نمی‌شود
		{Title: "اسلایدر پیش‌نویس (نمونه)", Slug: demoSliderSlugPrefix + "draft",
			Position: entities.SliderPositionBeforeFooter,
			Status:   entities.ProductStatusDraft, Products: 3, PreferDemo: true},
	}
}

func seedDemoSliders(db *gorm.DB) {
	var count int64
	db.Model(&entities.ProductSlider{}).Where("slug LIKE ?", demoSliderSlugPrefix+"%").Count(&count)
	if count > 0 {
		fmt.Printf("[seed] product sliders ....... already seeded (%d)\n", count)
		return
	}

	now := time.Now()
	for _, spec := range demoSliders(now) {
		startsAt, endsAt := spec.StartsAt, spec.EndsAt

		slider := entities.ProductSlider{
			Title:    spec.Title,
			Slug:     spec.Slug,
			Position: spec.Position,
			Status:   spec.Status,
			StartsAt: startsAt,
			EndsAt:   endsAt,
		}
		if catID := demoCategoryID(db, spec.CategorySlug); catID > 0 {
			slider.CategoryID = &catID
		}
		if !entities.ValidSliderPosition(slider.Position) {
			slider.Position = entities.SliderPositionAfterCategories
		}
		if err := db.Create(&slider).Error; err != nil {
			fmt.Printf("[seed] slider %q failed: %v\n", spec.Title, err)
			continue
		}

		for i, productID := range demoSliderProductIDs(db, spec.Products, spec.PreferDemo) {
			link := entities.SliderProduct{
				SliderID:  slider.ID,
				ProductID: productID,
				SortOrder: i,
			}
			if err := db.Create(&link).Error; err != nil {
				fmt.Printf("[seed]   slider product %d failed: %v\n", productID, err)
			}
		}
	}

	fmt.Println("[seed] product sliders ....... done (4 positions + draft, ≤10 products each)")
}

// ------------------------------------------------------------
// helpers
// ------------------------------------------------------------

// demoSliderProductIDs picks published products for a slider, preferring the
// demo set (they carry images) and capping at entities.MaxSliderProducts.
func demoSliderProductIDs(db *gorm.DB, want int, preferDemo bool) []uint {
	if want > entities.MaxSliderProducts {
		want = entities.MaxSliderProducts
	}

	order := "id ASC"
	if preferDemo {
		order = "CASE WHEN sku LIKE '" + demoPrefix + "%' THEN 0 ELSE 1 END, id ASC"
	}

	var ids []uint
	db.Model(&entities.Product{}).
		Where("status = ?", entities.ProductStatusPublished).
		Order(order).
		Limit(want).
		Pluck("id", &ids)
	return ids
}

// demoAttrValueID resolves "attributeCode:value" (e.g. "size:m") to the
// attribute_values id. Returns 0 when the value does not exist.
func demoAttrValueID(db *gorm.DB, pair string) uint {
	parts := strings.SplitN(pair, ":", 2)
	if len(parts) != 2 {
		return 0
	}
	code := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var id uint
	db.Raw(`
		SELECT av.id
		FROM attribute_values av
		         JOIN attributes a ON a.id = av.attribute_id
		WHERE a.code = ? AND av.value = ?
		  AND av.deleted_at IS NULL AND a.deleted_at IS NULL
		LIMIT 1`, code, value).Scan(&id)
	return id
}

func demoCategoryID(db *gorm.DB, slug string) uint {
	if slug == "" {
		return 0
	}
	var id uint
	db.Raw(`SELECT id FROM categories WHERE slug = ? AND deleted_at IS NULL LIMIT 1`, slug).Scan(&id)
	return id
}

func demoBrandID(db *gorm.DB, slug string) uint {
	if slug == "" {
		return 0
	}
	var id uint
	db.Raw(`SELECT id FROM brands WHERE slug = ? AND deleted_at IS NULL LIMIT 1`, slug).Scan(&id)
	return id
}

// copyBannerAsset copies a theme banner into the upload folder so the seeded
// banners resolve to a real image URL (returns the relative stored path).
func copyBannerAsset(name string) (string, error) {
	src := filepath.Join(bannerAssetDir, name)
	dst := filepath.Join(bannerUploadDir, name)

	if _, err := os.Stat(dst); err == nil {
		return filepath.ToSlash(filepath.Join("seed", name)), nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}

	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Join("seed", name)), nil
}

// ------------------------------------------------------------
// 4. order fee tariffs (هزینه ارسال / هزینه بسته‌بندی)
//
// Baseline rates so checkout quotes work out of the box:
//   - shipping 104,000 تومان, free above 1,000,000 تومان of items
//   - packaging 23,000 تومان, no threshold
// ------------------------------------------------------------

func seedFeeRates(db *gorm.DB) {
	specs := []struct {
		kind      string
		title     string
		amount    uint
		threshold *uint
	}{
		{entities.FeeKindShipping, "تعرفه ارسال پیش‌فرض", 104_000, uintPtrSeed(1_000_000)},
		{entities.FeeKindPackaging, "تعرفه بسته‌بندی پیش‌فرض", 23_000, nil},
	}

	for _, spec := range specs {
		var count int64
		db.Model(&entities.FeeRate{}).
			Where("kind = ? AND title = ?", spec.kind, spec.title).
			Count(&count)
		if count > 0 {
			fmt.Printf("[seed] fee rate %-28q already seeded\n", spec.title)
			continue
		}
		rate := entities.FeeRate{
			Kind:          spec.kind,
			Title:         spec.title,
			Amount:        spec.amount,
			FreeThreshold: spec.threshold,
			Status:        true,
		}
		// Select forces the zero-ok columns the same way the admin panel does
		if err := db.
			Select("Kind", "Title", "Amount", "FreeThreshold", "StartsAt", "EndsAt", "Status").
			Create(&rate).Error; err != nil {
			fmt.Printf("[seed] fee rate %q failed: %v\n", spec.title, err)
			continue
		}
	}

	fmt.Println("[seed] fee rates ............. done (shipping 104,000 / free ≥ 1,000,000 · packaging 23,000)")
}

func uintPtrSeed(v uint) *uint { return &v }
