package product

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
	"math"
	"shop/internal/database/mongodb"
	"shop/internal/entities"
	"shop/internal/pkg/util"
)

// SyncMongo وظیفه ذخیره محصول و روابطش و نیز ذخیره موجودی انبار رو در استراکت دلخواه در مونگو دیبی به عهده داره
func SyncMongo(c context.Context, db *gorm.DB, productID uint) error {

	// بارگذاری اطلاعات محصول
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
		fmt.Println("--- mongo product error:", productErr)
		return productErr
	}

	// تعریف ساختار مورد نیاز برای انبارها و ویژگی‌ها
	type InventoryWithAttributes struct {
		InventoryID                 uint
		Quantity                    uint
		AttributeID                 uint
		AttributeTitle              string
		AttributeValueID            uint
		AttributeValueTitle         string
		ProductInventoryAttributeID uint
	}

	// بارگذاری موجودی‌ها
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
		log.Println("--- mongo scan error:", serr)
		return serr
	}

	// آماده‌سازی برای ذخیره در MongoDB
	inventoryMap := make(map[string]entities.Inventory) // prepare product inventory

	for _, inventory := range inventories {
		key := fmt.Sprintf("%d", inventory.InventoryID) //convert int to string
		if _, exists := inventoryMap[key]; !exists {
			inventoryMap[key] = entities.Inventory{
				InventoryID: int64(inventory.InventoryID),
				Quantity:    int64(inventory.Quantity),
				Attributes:  []entities.InventoryAttributes{},
			}
		}

		//بعضی محصولات attribute ندارند و فقط موجودی دارند
		// prepare attributes array for inventories field (note : some products has no any attributes they just have inventory_id and quantity )
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

	// تبدیل محصول به ساختار MongoProduct
	mongoProduct := entities.MongoProduct{
		Product: entities.P{
			ID: int64(product.ID),
			Category: entities.C{
				ID:       int64(product.Category.ID),
				ParentID: int64(*product.Category.ParentID),
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
			Discount: func() int64 {
				originalPrice := float64(product.OriginalPrice)
				salePrice := float64(product.SalePrice)
				dis := ((originalPrice - salePrice) / originalPrice) * 100

				return int64(math.Round(dis))
			}(),
			Description: product.Description,
			Images:      entities.Img{Data: transformImages(product.ProductImages)},
			Features:    entities.F{Data: transformFeatures(product.Features)},
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
		Inventories: inventoryMap,
	}

	// چک کردن وجود محصول در MongoDB
	productsCollection := mongodb.GetCollection(mongodb.ProductsCollection)
	filter := bson.M{"product.id": mongoProduct.Product.ID}
	update := bson.M{"$set": mongoProduct}

	// تلاش برای آپدیت محصول در صورت وجود، در غیر این صورت ایجاد محصول جدید
	opts := options.Update().SetUpsert(true)
	_, err := productsCollection.UpdateOne(c, filter, update, opts)
	if err != nil {
		log.Println("--- product insert/update err ", err)
		return err
	}

	// upsert product in typesence
	go func() {
		util.UpsertInTypesence(c, util.UpsertTypesenceProduct{
			ID:    fmt.Sprintf("%d", product.ID),
			Title: product.Title,
			Slug:  product.Slug,
			Sku:   product.Sku,
		})
	}()

	log.Println("-- upsert product in mongoDB successfully")
	return nil
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
