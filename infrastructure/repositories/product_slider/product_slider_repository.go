package slider

import (
	"context"
	"fmt"
	"strings"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	responses "shop/application/dto/admin"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"shop/pkg/pagination"
)

type ProductSliderRepository struct {
	db *gorm.DB
}

func NewProductSliderRepository(db *gorm.DB) repositories.ProductSliderRepositoryInterface {
	return &ProductSliderRepository{db: db}
}

func (r *ProductSliderRepository) GetAll(c *gin.Context) ([]*entities.ProductSlider, error) {
	var sliders []*entities.ProductSlider
	// Products + Products.Product feed the list preview (count + first titles)
	err := r.db.WithContext(c).
		Preload("Category").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Preload("Products.Product").
		Order("position ASC, id DESC").
		Find(&sliders).Error
	return sliders, err
}

func (r *ProductSliderRepository) FindByID(c *gin.Context, sliderID uint) (*entities.ProductSlider, error) {
	return r.find(c, "id = ?", sliderID)
}

func (r *ProductSliderRepository) FindBySlug(c *gin.Context, slug string) (*entities.ProductSlider, error) {
	return r.find(c, "slug = ?", slug)
}

func (r *ProductSliderRepository) find(c context.Context, cond string, args ...interface{}) (*entities.ProductSlider, error) {
	var slider entities.ProductSlider
	err := r.db.WithContext(c).
		Preload("Category").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Preload("Products.Product").
		Preload("Products.Product.ProductImages").
		Where(cond, args...).
		First(&slider).Error
	if err != nil {
		return nil, err
	}
	return &slider, nil
}

// Store creates the slider and its ordered product set in one transaction.
func (r *ProductSliderRepository) Store(c *gin.Context, req *requests.CreateProductSliderRequest) (*entities.ProductSlider, error) {
	startsAt, endsAt := requests.SliderDates(req.StartsAt, req.EndsAt)
	slider := &entities.ProductSlider{
		Title:    strings.TrimSpace(req.Title),
		Slug:     provisionalSlug(req.Title),
		Position: req.Position,
		Status:   req.Status,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}
	if id := req.CatalogCategoryID(); id > 0 {
		slider.CategoryID = &id
	}

	txErr := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(slider).Error; err != nil {
			return err
		}
		// nice, unique URL: <slugified title>-<id> (Persian titles fall back
		// to "slider", so the id keeps the slug unique)
		base := slugify(slider.Title)
		if base == "" {
			base = "slider"
		}
		if err := tx.Model(slider).Update("slug", fmt.Sprintf("%s-%d", base, slider.ID)).Error; err != nil {
			return err
		}
		return r.replaceProducts(tx, slider.ID, req.ProductIDs)
	})
	if txErr != nil {
		return nil, txErr
	}
	return r.FindByID(c, slider.ID)
}

// Update rewrites the slider columns and replaces its product set.
func (r *ProductSliderRepository) Update(c *gin.Context, sliderID uint, req *requests.CreateProductSliderRequest) error {
	startsAt, endsAt := requests.SliderDates(req.StartsAt, req.EndsAt)
	values := map[string]interface{}{
		"title":     strings.TrimSpace(req.Title),
		"position":  req.Position,
		"status":    req.Status,
		"starts_at": startsAt,
		"ends_at":   endsAt,
	}
	if id := req.CatalogCategoryID(); id > 0 {
		values["category_id"] = id
	} else {
		values["category_id"] = nil
	}

	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.ProductSlider{}).Where("id = ?", sliderID).Updates(values).Error; err != nil {
			return err
		}
		return r.replaceProducts(tx, sliderID, req.ProductIDs)
	})
}

// Delete soft-deletes the slider and its product links.
func (r *ProductSliderRepository) Delete(c *gin.Context, sliderID uint) error {
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("slider_id = ?", sliderID).Delete(&entities.SliderProduct{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entities.ProductSlider{}, sliderID).Error
	})
}

// replaceProducts hard-deletes the current links and inserts the new set
// (keeps UNIQUE(slider_id, product_id) happy across edits).
func (r *ProductSliderRepository) replaceProducts(tx *gorm.DB, sliderID uint, productIDs []uint) error {
	if err := tx.Unscoped().Where("slider_id = ?", sliderID).Delete(&entities.SliderProduct{}).Error; err != nil {
		return err
	}
	for i, productID := range productIDs {
		link := entities.SliderProduct{
			SliderID:  sliderID,
			ProductID: productID,
			SortOrder: i,
		}
		if err := tx.Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}

// ActiveSliders returns the sliders the storefront may render right now.
func (r *ProductSliderRepository) ActiveSliders(ctx context.Context) ([]*entities.ProductSlider, error) {
	var sliders []*entities.ProductSlider
	err := r.db.WithContext(ctx).
		Preload("Category").
		// soft-delete scope already covers slider_products; order + cap
		// the curated set at MaxSliderProducts
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC").
				Limit(entities.MaxSliderProducts)
		}).
		Preload("Products.Product").
		Preload("Products.Product.ProductImages").
		Where("status = ?", entities.ProductStatusPublished).
		Where("(starts_at IS NULL OR starts_at <= CURRENT_DATE)").
		Where("(ends_at IS NULL OR ends_at >= CURRENT_DATE)").
		Order("position ASC, id ASC").
		Find(&sliders).Error
	if err != nil {
		return nil, err
	}
	return sliders, nil
}

// Catalog builds the «مشاهده همه» page:
//   - slider scoped to a category → every published product of that category
//     (the curated set is a subset of it)
//   - otherwise → the slider's own products
func (r *ProductSliderRepository) Catalog(c *gin.Context, slider *entities.ProductSlider) (pagination.Pagination, error) {
	pg := pagination.Pagination{
		Limit: pageLimit(c),
		Page:  pageNumber(c),
	}

	// curated slider without a category scope: the whole set fits in one page
	if slider.CategoryID == nil {
		products := make([]*entities.Product, 0, len(slider.Products))
		for _, link := range slider.Products {
			if link.Product != nil {
				products = append(products, link.Product)
			}
		}
		pg.TotalRows = int64(len(products))
		pg.TotalPages = 1
		pg.Page = 1
		pg.CurrentLink = "?page="
		pg.Rows = responses.ToProducts(products)
		return pg, nil
	}

	condition := fmt.Sprintf("category_id = %d AND status = '%s'",
		*slider.CategoryID, entities.ProductStatusPublished)

	var products []*entities.Product
	paginateQuery, exist := pagination.Paginate(c, condition, &products, &pg, r.db)
	if !exist {
		return pg, gorm.ErrRecordNotFound
	}
	if err := paginateQuery(r.db).
		Preload("Category").
		Preload("ProductImages").
		Where("category_id = ? AND status = ?", *slider.CategoryID, entities.ProductStatusPublished).
		Order("id DESC").
		Find(&products).Error; err != nil {
		return pg, err
	}

	pg.Rows = responses.ToProducts(products)
	return pg, nil
}

// --- small helpers -------------------------------------------------------

func pageLimit(c *gin.Context) int {
	n, err := strconv.Atoi(c.Query("limit"))
	if err != nil || n < 1 {
		return 12
	}
	if n > 48 {
		return 48
	}
	return n
}

func pageNumber(c *gin.Context) int {
	n, err := strconv.Atoi(c.Query("page"))
	if err != nil || n < 1 {
		return 1
	}
	return n
}

// provisionalSlug gives the row a unique slug before the real one is written.
func provisionalSlug(title string) string {
	return fmt.Sprintf("slider-%d", time.Now().UnixNano())
}

// slugify mirrors the product slug helper (ASCII only, Persian → empty).
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
