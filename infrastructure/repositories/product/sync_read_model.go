package product

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/pkg/util"
)

// SyncReadModel builds the flattened product document (product + category + brand
// + images + features + inventories with attributes) and persists it into
// products.read_model (JSONB). It replaces the old MongoDB-based SyncMongo.
// After the JSONB write it upserts the product into Typesense.
func SyncReadModel(c context.Context, db *gorm.DB, productID uint) error {
	var product entities.Product
	productErr := db.WithContext(c).
		Preload("Category").
		Preload("Brand").
		Preload("ProductImages").
		Preload("Features").
		Where("id=?", productID).
		First(&product).
		Error
	if productErr != nil {
		fmt.Println("--- read model build error:", productErr)
		return productErr
	}

	type InventoryWithAttributes struct {
		InventoryID                 uint
		Quantity                    uint
		ReservedStock               uint
		Price                       *uint
		SalePrice                   *uint
		DiscountPrice               *uint
		Status                      string
		AttributeID                 uint
		AttributeTitle              string
		AttributeValueID            uint
		AttributeValueTitle         string
		ProductInventoryAttributeID uint
	}

	var inventories []InventoryWithAttributes
	serr := db.
		WithContext(c).
		Table("product_variants").
		Select("product_variants.id AS inventory_id, product_variants.stock AS quantity, product_variants.reserved_stock, product_variants.price, product_variants.sale_price, product_variants.discount_price, product_variants.status, attributes.id AS attribute_id, attributes.title AS attribute_title, attribute_values.id AS attribute_value_id, attribute_values.value AS attribute_value_title, variant_attribute_values.id AS product_inventory_attribute_id").
		Joins("LEFT JOIN variant_attribute_values ON product_variants.id = variant_attribute_values.variant_id AND variant_attribute_values.deleted_at IS NULL").
		Joins("LEFT JOIN attribute_values ON variant_attribute_values.attribute_value_id = attribute_values.id AND attribute_values.deleted_at IS NULL").
		Joins("LEFT JOIN attributes ON attribute_values.attribute_id = attributes.id AND attributes.deleted_at IS NULL").
		Where("product_variants.product_id = ? AND product_variants.deleted_at IS NULL", product.ID).
		Scan(&inventories).
		Error

	if serr != nil {
		log.Println("--- read model inventory scan error:", serr)
		return serr
	}

	inventoryMap := make(map[string]entities.Inventory)
	for _, inventory := range inventories {
		key := fmt.Sprintf("%d", inventory.InventoryID)
		if _, exists := inventoryMap[key]; !exists {
			inventoryMap[key] = variantInventory(inventory, product)
		}
		inv := inventoryMap[key]
		inv.Attributes = append(inv.Attributes, entities.InventoryAttributes{
			AttributeID:                 int64(inventory.AttributeID),
			AttributeTitle:              inventory.AttributeTitle,
			AttributeValueID:            int64(inventory.AttributeValueID),
			AttributeValueTitle:         inventory.AttributeValueTitle,
			ProductInventoryAttributeID: int64(inventory.ProductInventoryAttributeID),
		})
		inventoryMap[key] = inv
	}

	// golden rule: the storefront discount is measured against the price the
	// customer actually pays (aggregate minimum when variants exist), not the
	// nominal product sale price.
	effectiveSale := int64(product.SalePrice)
	if product.MinPrice > 0 {
		effectiveSale = int64(product.MinPrice)
	}
	discount := int64(0)
	if product.OriginalPrice > 0 && effectiveSale > 0 && effectiveSale < int64(product.OriginalPrice) {
		discount = int64(math.Round((float64(product.OriginalPrice) - float64(effectiveSale)) / float64(product.OriginalPrice) * 100))
	}

	readModel := entities.ProductReadModel{
		Product: entities.P{
			ID: int64(product.ID),
			Category: entities.C{
				ID:       int64(product.Category.ID),
				ParentID: derefUintToInt64(product.Category.ParentID),
				Title:    product.Category.Title,
				Slug:     product.Category.Slug,
			},
			CategoryID: int64(product.CategoryID),
			Brand: entities.B{
				ID:    int64(product.Brand.ID),
				Title: product.Brand.Title,
				Slug:  product.Brand.Slug,
			},
			BrandID:       int64(product.BrandID),
			Title:         product.Title,
			Slug:          product.Slug,
			Sku:           product.Sku,
			Status:        product.Status,
			OriginalPrice: int64(product.OriginalPrice),
			SalePrice:     int64(product.SalePrice),
			Discount:      discount,
			MinPrice:      int64(product.MinPrice),
			MaxPrice:      int64(product.MaxPrice),
			Description:   product.Description,
			Images:        entities.Img{Data: transformImages(product.ProductImages)},
			Features:      entities.F{Data: transformFeatures(product.Features)},
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
		},
		Inventories: inventoryMap,
	}

	doc, err := json.Marshal(readModel)
	if err != nil {
		log.Println("--- read model marshal error:", err)
		return err
	}

	if err := db.WithContext(c).
		Model(&entities.Product{}).
		Where("id = ?", productID).
		Update("read_model", doc).
		Error; err != nil {
		log.Println("--- update read_model error:", err)
		return err
	}

	// sync the rich document into typesense: upsert while published,
	// delete while draft/archived (realtime search stays consistent)
	idStr := fmt.Sprintf("%d", product.ID)
	if product.Status != entities.ProductStatusPublished {
		go util.DeleteInTypesence(c, idStr)
	} else {
		// prefer the PricingService aggregates (Laravel min_price/in_stock);
		// fall back to the variant rows when the cache was not refreshed yet
		var totalStock int64
		for _, inv := range inventoryMap {
			totalStock += inv.Quantity
		}
		if product.TotalStock > 0 {
			totalStock = int64(product.TotalStock)
		}
		inStock := totalStock-int64(product.TotalReserved) > 0
		if product.AvailableStock > 0 || product.TotalStock > 0 {
			inStock = product.InStock
		}

		// effective selling price = min effective variant price (golden rule #5)
		salePrice := product.SalePrice
		if product.MinPrice > 0 {
			salePrice = product.MinPrice
		}
		// discount % against the effective selling price
		saleDiscount := discount
		if product.OriginalPrice > 0 && salePrice != product.SalePrice {
			pct := (float64(product.OriginalPrice) - float64(salePrice)) / float64(product.OriginalPrice) * 100
			if pct < 0 {
				pct = 0
			}
			saleDiscount = int64(math.Round(pct))
		}

		categoryTitle := ""
		if product.Category != nil {
			categoryTitle = product.Category.Title
		}
		brandTitle := ""
		if product.Brand != nil {
			brandTitle = product.Brand.Title
		}
		go util.UpsertInTypesence(c, util.UpsertTypesenceProduct{
			ID:            idStr,
			Title:         product.Title,
			Slug:          product.Slug,
			Sku:           product.Sku,
			Description:   product.Description,
			Category:      categoryTitle,
			Brand:         brandTitle,
			OriginalPrice: int64(product.OriginalPrice),
			SalePrice:     int64(salePrice),
			Discount:      int64(saleDiscount),
			Stock:         totalStock,
			InStock:       inStock,
			Status:        product.Status,
		})
	}

	log.Println("-- update product read_model (jsonb) successfully, product id:", productID)
	return nil
}

// variantInventory resolves one variant's pricing/status for the read model:
// NULL price columns inherit the product price, discount counts only when it
// is greater than zero and lower than the sale price (golden rule #3).
func variantInventory(row struct {
	InventoryID                 uint
	Quantity                    uint
	ReservedStock               uint
	Price                       *uint
	SalePrice                   *uint
	DiscountPrice               *uint
	Status                      string
	AttributeID                 uint
	AttributeTitle              string
	AttributeValueID            uint
	AttributeValueTitle         string
	ProductInventoryAttributeID uint
}, product entities.Product) entities.Inventory {
	price := int64(product.OriginalPrice)
	if row.Price != nil {
		price = int64(*row.Price)
	}
	sale := int64(product.SalePrice)
	if row.SalePrice != nil {
		sale = int64(*row.SalePrice)
	}

	effective := sale
	hasDiscount := false
	// the badge percent is measured against the crossed-out list price (price),
	// consistent with the storefront display next to it
	var discountPercent int64
	if row.DiscountPrice != nil && *row.DiscountPrice > 0 && int64(*row.DiscountPrice) < sale {
		effective = int64(*row.DiscountPrice)
		hasDiscount = true
		if price > 0 && effective < price {
			discountPercent = int64(math.Round(float64(price-effective) / float64(price) * 100))
		}
	}

	available := int64(0)
	if row.Quantity > row.ReservedStock {
		available = int64(row.Quantity - row.ReservedStock)
	}

	return entities.Inventory{
		InventoryID:     int64(row.InventoryID),
		Quantity:        int64(row.Quantity),
		Attributes:      []entities.InventoryAttributes{},
		Price:           price,
		SalePrice:       sale,
		DiscountPrice:   int64(derefUint(row.DiscountPrice)),
		HasDiscount:     hasDiscount,
		DiscountPercent: discountPercent,
		EffectivePrice:  effective,
		Status:          row.Status,
		Available:       available,
	}
}

func derefUint(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

func derefUintToInt64(p *uint) int64 {
	if p == nil {
		return 0
	}
	return int64(*p)
}

func transformImages(images []*entities.ProductImages) []entities.ImgData {
	var imgData []entities.ImgData
	for _, img := range images {
		imgData = append(imgData, entities.ImgData{
			ID:           int64(img.ID),
			OriginalPath: img.Path,
			FullPath:     viper.GetString("Upload.Products") + img.Path,
		})
	}
	return imgData
}

func transformFeatures(features []*entities.Feature) []entities.FData {
	var fData []entities.FData
	for _, feature := range features {
		fData = append(fData, entities.FData{
			ID:        int64(feature.ID),
			ProductID: int64(feature.ProductID),
			Title:     feature.Title,
			Value:     feature.Value,
		})
	}
	return fData
}
