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
		AttributeID                 uint
		AttributeTitle              string
		AttributeValueID            uint
		AttributeValueTitle         string
		ProductInventoryAttributeID uint
	}

	var inventories []InventoryWithAttributes
	serr := db.
		WithContext(c).
		Table("product_inventories").
		Select("product_inventories.id AS inventory_id, product_inventories.quantity, product_attributes.attribute_id, attributes.title AS attribute_title, attribute_values.id AS attribute_value_id, attribute_values.value AS attribute_value_title, product_inventory_attributes.id AS product_inventory_attribute_id").
		Joins("LEFT JOIN product_inventory_attributes ON product_inventories.id = product_inventory_attributes.product_inventory_id AND product_inventory_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN product_attributes ON product_inventory_attributes.product_attribute_id = product_attributes.id AND product_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attributes ON product_attributes.attribute_id = attributes.id AND attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attribute_values ON product_attributes.attribute_value_id = attribute_values.id AND attribute_values.deleted_at IS NULL").
		Where("product_inventories.product_id = ? AND product_inventories.deleted_at IS NULL", product.ID).
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
			inventoryMap[key] = entities.Inventory{
				InventoryID: int64(inventory.InventoryID),
				Quantity:    int64(inventory.Quantity),
				Attributes:  []entities.InventoryAttributes{},
			}
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

	discount := int64(0)
	if product.OriginalPrice > 0 {
		originalPrice := float64(product.OriginalPrice)
		salePrice := float64(product.SalePrice)
		discount = int64(math.Round(((originalPrice - salePrice) / originalPrice) * 100))
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

	// upsert product in typesense
	go util.UpsertInTypesence(c, util.UpsertTypesenceProduct{
		ID:    fmt.Sprintf("%d", product.ID),
		Title: product.Title,
		Slug:  product.Slug,
		Sku:   product.Sku,
	})

	log.Println("-- update product read_model (jsonb) successfully, product id:", productID)
	return nil
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
