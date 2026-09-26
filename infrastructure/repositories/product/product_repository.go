package product

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepositoryInterface {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]*entities.Product, error) {
	var products []*entities.Product
	err := p.db.WithContext(ctx).Preload("Category").Preload("Brand").Find(&products).Error
	return products, err
}

func (p *ProductRepository) FindBy(ctx context.Context, columnName string, value any) (*entities.Product, error) {
	var product entities.Product
	condition := fmt.Sprintf("%s=?", columnName)
	err := p.db.
		WithContext(ctx).
		Preload("Category").Preload("Brand").
		Preload("ProductAttributes").Preload("ProductVariants").
		Preload("ProductVariants.VariantAttributeValues.AttributeValue").
		Preload("VariantAttributeValues.AttributeValue").
		Preload("ProductImages").Preload("Features").
		First(&product, condition, value).
		Error

	return &product, err
}

func (p *ProductRepository) FindByID(ctx context.Context, ID int) (*entities.Product, error) {
	var product entities.Product
	err := p.db.WithContext(ctx).Preload("ProductAttributes").First(&product, ID).Error
	return &product, err
}

// GetList builds the admin product list query (Laravel ProductController::index):
// text search, status / category / in_stock filters, attributes_json facets and
// sorting on the aggregate cache columns.
func (p *ProductRepository) GetList(ctx context.Context, q requests.ProductListQuery) ([]*entities.Product, error) {
	var products []*entities.Product

	db := p.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("ProductImages")

	if search := strings.TrimSpace(q.Q); search != "" {
		like := "%" + escapeLike(search) + "%"
		db = db.Where("(title ILIKE ? OR sku ILIKE ?)", like, like)
	}
	if status := strings.TrimSpace(q.Status); status != "" {
		db = db.Where("status = ?", status)
	}
	if q.CategoryID > 0 {
		db = db.Where("category_id = ?", q.CategoryID)
	}
	if q.BrandID > 0 {
		db = db.Where("brand_id = ?", q.BrandID)
	}
	if q.InStock {
		db = db.Where("in_stock = ?", true)
	}
	for code, valueIDs := range q.Attrs {
		if code == "" || len(valueIDs) == 0 {
			continue
		}
		conds := make([]string, 0, len(valueIDs))
		args := make([]interface{}, 0, len(valueIDs)*2)
		for _, id := range valueIDs {
			conds = append(conds, "(attributes_json -> ?) @> ?::jsonb")
			args = append(args, code, fmt.Sprintf("[%d]", id))
		}
		db = db.Where("("+strings.Join(conds, " OR ")+")", args...)
	}

	switch q.Sort {
	case "price_asc":
		db = db.Order("min_price ASC, id DESC")
	case "price_desc":
		db = db.Order("max_price DESC, id DESC")
	case "name":
		db = db.Order("title ASC")
	default:
		db = db.Order("created_at DESC")
	}

	err := db.Find(&products).Error
	return products, err
}

// escapeLike neutralizes LIKE wildcards coming from user input.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "\\\\", "\\\\\\\\")
	s = strings.ReplaceAll(s, "%", "\\\\%")
	s = strings.ReplaceAll(s, "_", "\\\\_")
	return s
}

// Store creates the product together with its variant rows (Laravel
// ProductController::store pattern): product + variants + attribute links in
// one transaction. The caller refreshes pricing aggregates and then the read
// model (see application/usecases/product.Create).
func (p *ProductRepository) Store(ctx context.Context, product *entities.Product, rows []requests.VariantRow) (*entities.Product, error) {

	txErr := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// no variant rows posted → one stock-only variant so the storefront
		// always has something sellable
		if len(rows) == 0 {
			rows = []requests.VariantRow{{
				Price:     product.OriginalPrice,
				SalePrice: product.SalePrice,
				Status:    entities.VariantStatusActive,
			}}
		}

		for _, row := range rows {
			if err := createVariantWithLinks(tx, product.ID, row); err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}
	return product, nil
}

// createVariantWithLinks inserts one variant and its variant_attribute_values
// links. Runs inside the caller's transaction.
func createVariantWithLinks(tx *gorm.DB, productID uint, row requests.VariantRow) error {
	variant := entities.ProductVariant{
		ProductID:     productID,
		Sku:           variantSku(row.Sku),
		Price:         &row.Price,
		SalePrice:     &row.SalePrice,
		DiscountPrice: row.DiscountPrice,
		Stock:         row.Stock,
		Status:        variantStatus(row.Status),
		ExpiresAt:     parseVariantExpiresAt(row.ExpiresAt),
	}
	if err := tx.Create(&variant).Error; err != nil {
		return err
	}
	return createVariantLinks(tx, productID, variant.ID, row.AttributeValueIDs)
}

func createVariantLinks(tx *gorm.DB, productID, variantID uint, valueIDs []uint) error {
	for _, valueID := range valueIDs {
		link := entities.VariantAttributeValue{
			ProductID:        productID,
			VariantID:        variantID,
			AttributeValueID: valueID,
		}
		if err := tx.Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}

func variantSku(sku string) *string {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil
	}
	return &sku
}

func variantStatus(status string) string {
	switch status {
	case entities.VariantStatusActive, entities.VariantStatusInactive:
		return status
	}
	return entities.VariantStatusActive
}

// updateVariantColumns writes the editable variant fields posted by the form.
func updateVariantColumns(tx *gorm.DB, variant *entities.ProductVariant, row requests.VariantRow) error {
	return tx.Model(variant).Updates(map[string]interface{}{
		"price":          row.Price,
		"sale_price":     row.SalePrice,
		"discount_price": row.DiscountPrice,
		"stock":          row.Stock,
		"status":         variantStatus(row.Status),
		"expires_at":     parseVariantExpiresAt(row.ExpiresAt),
	}).Error
}

// Update writes the product columns and its variant rows in one transaction:
// rows carrying an id update that variant, rows without an id are new
// combinations (product_type = variable). For product_type = simple the single
// variant is updated from the base price rows.
func (p *ProductRepository) Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) (*entities.Product, error) {

	var product entities.Product
	pErr := p.db.WithContext(c).First(&product, productID).Error
	if pErr != nil {
		fmt.Println("---- repo product find err : ", pErr)
		return nil, pErr
	}

	status := strings.TrimSpace(req.Status)
	if !entities.ValidProductStatus(status) {
		status = product.Status
	}

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&product).Updates(map[string]interface{}{
			"category_id":    req.CategoryID,
			"brand_id":       req.BrandID,
			"title":          strings.TrimSpace(req.Title),
			"slug":           strings.TrimSpace(req.Slug),
			"sku":            strings.TrimSpace(req.Sku),
			"status":         status,
			"expires_at":     req.ExpiresAtTime(),
			"original_price": req.OriginalPrice,
			"sale_price":     req.SalePrice,
			"description":    req.Description,
		}).Error; err != nil {
			return err
		}

		if len(req.Variants) > 0 {
			// variable product: update by id, create the rest (a posted row
			// whose combination already exists updates it instead of
			// duplicating it)
			combos, err := existingVariantCombos(tx, uint(productID))
			if err != nil {
				return err
			}
			for _, row := range req.Variants {
				if row.ID > 0 {
					var variant entities.ProductVariant
					if err := tx.Where("id = ? AND product_id = ?", row.ID, productID).
						First(&variant).Error; err != nil {
						return err
					}
					if err := updateVariantColumns(tx, &variant, row); err != nil {
						return err
					}
					continue
				}
				if id, ok := combos[valueIDsKey(row.AttributeValueIDs)]; ok && id > 0 {
					var variant entities.ProductVariant
					if err := tx.First(&variant, id).Error; err != nil {
						return err
					}
					if err := updateVariantColumns(tx, &variant, row); err != nil {
						return err
					}
					continue
				}
				if err := createVariantWithLinks(tx, uint(productID), row); err != nil {
					return err
				}
			}
			return nil
		}

		// simple product: keep a single variant in sync with the base rows
		row := req.SimpleVariantRows()[0]
		var variant entities.ProductVariant
		err := tx.Where("product_id = ?", productID).Order("id").First(&variant).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return createVariantWithLinks(tx, uint(productID), row)
		}
		return updateVariantColumns(tx, &variant, row)
	})

	if txErr != nil {
		return nil, txErr
	}

	return &product, nil
}

// valueIDsKey builds a stable key for one combination (sorted value ids).
func valueIDsKey(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	cp := append([]uint(nil), ids...)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	parts := make([]string, 0, len(cp))
	for _, id := range cp {
		parts = append(parts, strconv.FormatUint(uint64(id), 10))
	}
	return strings.Join(parts, ",")
}

// existingVariantCombos maps combination keys → variant id for one product so
// an update never creates a duplicate combination.
func existingVariantCombos(tx *gorm.DB, productID uint) (map[string]uint, error) {
	type row struct {
		VariantID uint
		ValueID   uint
	}
	var rows []row
	if err := tx.Table("variant_attribute_values").
		Select("variant_attribute_values.variant_id AS variant_id, variant_attribute_values.attribute_value_id AS value_id").
		Joins("JOIN product_variants v ON v.id = variant_attribute_values.variant_id AND v.deleted_at IS NULL").
		Where("variant_attribute_values.product_id = ? AND variant_attribute_values.deleted_at IS NULL", productID).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	byVariant := map[uint][]uint{}
	for _, r := range rows {
		byVariant[r.VariantID] = append(byVariant[r.VariantID], r.ValueID)
	}
	out := map[string]uint{}
	for variantID, ids := range byVariant {
		out[valueIDsKey(ids)] = variantID
	}
	return out, nil
}

// parseVariantExpiresAt converts the posted YYYY-MM-DD value to a *time.Time
// (nil when empty or malformed — validation happens in the request layer).
func parseVariantExpiresAt(v string) *time.Time {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return &t
		}
	}
	return nil
}

// SyncReadModel rebuilds products.read_model and syncs Typesense for one
// product (ProductRepositoryInterface).
func (p *ProductRepository) SyncReadModel(ctx context.Context, productID uint) error {
	return SyncReadModel(ctx, p.db, productID)
}
