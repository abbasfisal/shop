package product

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"shop/internal/database/mongodb"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"strconv"
	"strings"
)

type ProductRepository struct {
	db          *gorm.DB
	mongoClient *mongo.Client
}

func NewProductRepository(db *gorm.DB, mongoClient *mongo.Client) ProductRepository {
	return ProductRepository{
		db:          db,
		mongoClient: mongoClient,
	}
}

func (p ProductRepository) GetAll(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	err := p.db.WithContext(ctx).Preload("Category").Preload("Brand").Find(&products).Error
	return products, err
}

func (p ProductRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Product, error) {
	var product entities.Product
	condition := fmt.Sprintf("%s=?", columnName)
	err := p.db.Preload("Category").Preload("Brand").Preload("ProductAttributes").Preload("ProductInventories").Preload("ProductImages").Preload("Features").First(&product, condition, value).Error
	return product, err
}

func (p ProductRepository) FindByID(ctx context.Context, ID int) (entities.Product, error) {
	var product entities.Product
	err := p.db.WithContext(ctx).Preload("ProductAttributes").First(&product, ID).Error
	return product, err
}

func (p ProductRepository) Store(ctx context.Context, product entities.Product) (entities.Product, error) {

	err := p.db.WithContext(ctx).Create(&product).Error

	if err == nil {
		productsCollection := mongodb.GetCollection(mongodb.ProductsCollection)

		_, err := productsCollection.InsertOne(ctx, &entities.MongoProduct{
			ID:            product.ID,
			CategoryID:    product.CategoryID,
			BrandID:       product.BrandID,
			Title:         product.Title,
			Slug:          product.Slug,
			Sku:           product.Sku,
			Status:        product.Status,
			OriginalPrice: product.OriginalPrice,
			SalePrice:     product.SalePrice,
			Description:   product.Description,
			Images: func() []string {
				var imgs []string
				for _, item := range product.ProductImages {
					imgs = append(imgs, item.Path)
				}
				return imgs
			}(),
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		})

		if err != nil {
			fmt.Println("\n --------- product mongodb : failed ---------- | err: ", err)
		}
		fmt.Println("\n------- product mongodb : created successful")
	}
	//for _, Ip := range product.ProductImage {
	//	fmt.Println("insed ", Ip, "| product id : ", product.ID)
	//	p.db.Create(entities.ProductImages{
	//		ProductID: product.ID,
	//		Path:      Ip.Path,
	//	})
	//}

	return product, err
}

func (p ProductRepository) GetRootAttributes(c *gin.Context, productID int) ([]entities.Attribute, error) {
	var category entities.Category
	var attributes []entities.Attribute

	var product entities.Product
	pErr := p.db.WithContext(c).Where("id = ? ", productID).First(&product).Error

	if pErr != nil {
		return attributes, pErr
	}

	err := p.db.Raw(
		` WITH RECURSIVE CategoryHierarchy AS (
            SELECT id, title, parent_id
            FROM categories
            WHERE id = ?

            UNION ALL

            SELECT c.id, c.title, c.parent_id
            FROM categories c
            INNER JOIN CategoryHierarchy ch ON c.id = ch.parent_id
        )
        SELECT *
        FROM CategoryHierarchy
        WHERE parent_id IS NULL
        LIMIT 1;`, product.CategoryID,
	).Scan(&category).Error

	if err != nil {
		fmt.Println("product repository _ root category not found")
		return attributes, err
	}

	aErr := p.db.WithContext(c).Preload("AttributeValues").Find(&attributes).Error

	return attributes, aErr
}

func (p ProductRepository) StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error {
	//find product by id
	_, err := p.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	//store []attributes values into product_attributes table
	for _, v := range attValues {
		parts := strings.Split(v, ":")

		attributeID, _ := strconv.Atoi(parts[1])
		attributeValueID, _ := strconv.Atoi(parts[5])
		p.db.Create(&entities.ProductAttribute{
			ProductID:           uint(productID),
			AttributeID:         uint(attributeID),
			AttributeTitle:      parts[3],
			AttributeValueID:    uint(attributeValueID),
			AttributeValueTitle: parts[7],
		})
	}
	return nil
}

func (p ProductRepository) GetProductAndAttributes(ctx *gin.Context, productID int) (map[string]interface{}, error) {
	type InventoryWithAttributes struct {
		InventoryID                 uint
		Quantity                    uint
		AttributeID                 uint
		AttributeTitle              string
		AttributeValueID            uint
		AttributeValueTitle         string
		ProductInventoryAttributeID uint
	}

	var product entities.Product
	aerr := p.db.WithContext(ctx).
		Where("id = ?", productID).
		First(&product).Error

	if aerr != nil {
		return map[string]interface{}{}, aerr
	}

	var inventories []InventoryWithAttributes

	result := make(map[string]interface{})

	serr := p.db.WithContext(ctx).
		Table("product_inventories").
		Select("product_inventories.id AS inventory_id, product_inventories.quantity, product_attributes.attribute_id, attributes.title AS attribute_title, attribute_values.id AS attribute_value_id, attribute_values.value AS attribute_value_title, product_inventory_attributes.id AS product_inventory_attribute_id").
		Joins("LEFT JOIN product_inventory_attributes ON product_inventories.id = product_inventory_attributes.product_inventory_id AND product_inventory_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN product_attributes ON product_inventory_attributes.product_attribute_id = product_attributes.id AND product_attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attributes ON product_attributes.attribute_id = attributes.id AND attributes.deleted_at IS NULL").
		Joins("LEFT JOIN attribute_values ON product_attributes.attribute_value_id = attribute_values.id AND attribute_values.deleted_at IS NULL").
		Where("product_inventories.product_id = ? and product_inventories.deleted_at IS NULL", productID).
		Scan(&inventories).Error

	if serr != nil {
		return map[string]interface{}{}, serr
	}

	inventoryMap := make(map[uint]map[string]interface{})
	for _, inventory := range inventories {
		if _, exists := inventoryMap[inventory.InventoryID]; !exists {
			inventoryMap[inventory.InventoryID] = map[string]interface{}{
				"add_attribute_link":    fmt.Sprintf("/admins/inventories/%d/attributes/add", inventory.InventoryID),  //add attribute-value to specific inventory
				"edit_inventory_link":   fmt.Sprintf("/admins/inventories/%d/update-quantity", inventory.InventoryID), //edit quantity of a product inventory (product_inventories)
				"delete_inventory_link": fmt.Sprintf("/admins/inventories/%d/delete", inventory.InventoryID),          //remove record from product_inventories table
				"quantity":              inventory.Quantity,
				"inventory_id":          inventory.InventoryID,
				"attributes":            []map[string]interface{}{},
			}
		}

		attributes := inventoryMap[inventory.InventoryID]["attributes"].([]map[string]interface{})
		attributes = append(attributes, map[string]interface{}{
			"attribute_id":                   inventory.AttributeID,
			"attribute_title":                inventory.AttributeTitle,
			"attribute_value_id":             inventory.AttributeValueID,
			"attribute_value_title":          inventory.AttributeValueTitle,
			"product_inventory_attribute_id": inventory.ProductInventoryAttributeID,
			"delete_attribute_link":          fmt.Sprintf("/admins/product-inventory-attributes/%d/delete", inventory.ProductInventoryAttributeID), //remove from product_inventory_attributes
		})
		inventoryMap[inventory.InventoryID]["attributes"] = attributes
	}

	result["product"] = product
	result["inventories"] = inventoryMap

	return result, nil
}

func (p ProductRepository) StoreProductInventory(c *gin.Context, productID int, req requests.CreateProductInventoryRequest) (entities.ProductInventory, error) {

	var inventory entities.ProductInventory

	//start transaction
	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var productAttributes []entities.ProductAttribute

		//fetch product-attributes
		//len(req.ProductAttributes)<=0  یعنی برای محصول ویژگی -مقدار نمیخواهیم بذاریم و صرفا میخواهیم موجودی بذاریم
		if len(req.ProductAttributes) > 0 {
			if err := tx.WithContext(c).Where("id IN ? ", req.ProductAttributes).Find(&productAttributes).Error; err != nil {
				return err
			}
			//check len retrieved product-attribute
			if len(productAttributes) != len(req.ProductAttributes) {
				return gorm.ErrRecordNotFound
			}
		}

		inventory = entities.ProductInventory{
			ProductID: uint(productID),
			Quantity:  uint(req.Quantity),
		}

		//todo: باید چک کنی که چندتا موجودی بدون اتریبیوت ذخیره شده تا بتونی روی ایجاد چندین موجودی بدون ویژگی کنترل داشته باشی
		//var count int64
		//if err := tx.Where("product_id = ?", productID).Count(&count).Error; err != nil {
		//	return err
		//}
		//if count > 1 {
		//	return &custom_error.DuplicateProductInventory{ProductID: uint(productID)}
		//}

		//store inventory
		if iErr := tx.WithContext(c).Create(&inventory).Error; iErr != nil {
			return iErr
		}

		//store product-attribute in product-inventory-attribute table
		//len(req.ProductAttributes)<=0  یعنی برای محصول ویژگی -مقدار نمیخواهیم بذاریم و صرفا میخواهیم موجودی بذاریم
		if len(req.ProductAttributes) > 0 {
			for _, attr := range productAttributes {
				inventoryAttr := entities.ProductInventoryAttribute{
					ProductID:          uint(productID),
					ProductInventoryID: inventory.ID,
					ProductAttributeID: attr.ID,
				}
				if err := tx.Create(&inventoryAttr).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if txErr != nil {
		fmt.Println("---- create inventory product err: ", txErr)
		return entities.ProductInventory{}, txErr
	}

	return inventory, nil
}

func (p ProductRepository) GetImage(c *gin.Context, imageID int) (entities.ProductImages, error) {
	var image entities.ProductImages
	err := p.db.Find(&image, imageID).Error
	return image, err
}

func (p ProductRepository) DeleteImage(c *gin.Context, imageID int) error {
	return p.db.WithContext(c).Unscoped().Delete(&entities.ProductImages{}, imageID).Error
}

func (p ProductRepository) StoreImages(c *gin.Context, productID int, imageStoredPath []string) error {
	var images []entities.ProductImages
	for _, image := range imageStoredPath {
		images = append(images, entities.ProductImages{
			ProductID: uint(productID),
			Path:      image,
		},
		)
	}

	return p.db.WithContext(c).Create(&images).Error
}

func (p ProductRepository) Update(c *gin.Context, productID int, req requests.UpdateProductRequest) (entities.Product, error) {

	var product entities.Product
	pErr := p.db.WithContext(c).First(&product, productID).Error
	if pErr != nil {
		fmt.Println("---- repo product find err : 182 ", pErr)
		return product, pErr
	}

	updateErr := p.db.WithContext(c).Model(&product).Update("category_id", req.CategoryID).
		Update("brand_id", req.BrandID).
		Update("title", strings.TrimSpace(req.Title)).
		Update("slug", strings.TrimSpace(req.Slug)).
		Update("sku", strings.TrimSpace(req.Sku)).
		Update("status", func() bool {
			if req.Status == "" {
				return false
			}
			return true
		}()).
		Update("original_price", req.OriginalPrice).
		Update("sale_price", req.SalePrice).
		Update("description", req.Description).Error

	if updateErr != nil {
		fmt.Println("---- repo product update  err : 201 ", updateErr)

		return entities.Product{}, pErr
	}
	fmt.Println("---- repo product udpate succ ")
	return product, nil
}

func (p ProductRepository) DeleteInventoryAttribute(c *gin.Context, productInventoryAttributeID int) error {

	//find
	var productInventoryAttribute entities.ProductInventoryAttribute
	if err := p.db.First(&productInventoryAttribute, productInventoryAttributeID).Error; err != nil {
		return err
	}

	//delete from product_inventory_attributes table
	if piaErr := p.db.WithContext(c).Unscoped().Delete(&productInventoryAttribute).Error; piaErr != nil {
		return piaErr
	}

	return nil
}

func (p ProductRepository) DeleteInventory(c *gin.Context, inventoryID int) error {

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		var inventory entities.ProductInventory

		//find inventory
		if iErr := p.db.WithContext(c).First(&inventory, inventoryID).Error; iErr != nil {
			return iErr
		}

		//delete all product-attribute inventory
		var productInventoryAttributes []entities.ProductInventoryAttribute
		if deleteErr := p.db.Where("product_inventory_id = ? ", inventory.ID).Delete(&productInventoryAttributes).Error; deleteErr != nil {
			return deleteErr
		}

		//delete inventory
		if iDelete := p.db.WithContext(c).Delete(&inventory).Error; iDelete != nil {
			return iDelete
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

func (p ProductRepository) AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) error {

	var productInventory entities.ProductInventory

	//find productInventory
	if err := p.db.WithContext(c).First(&productInventory, inventoryID).Error; err != nil {
		return err
	}

	//start transaction
	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var productAttributes []entities.ProductAttribute

		//fetch product-attributes
		if err := p.db.WithContext(c).Where("id IN ? ", attributes).Find(&productAttributes).Error; err != nil {
			return err
		}
		//check len retrieved product-attribute
		if len(productAttributes) != len(attributes) {
			return gorm.ErrRecordNotFound
		}

		//store product-attribute in product-inventory-attribute table
		for _, attr := range productAttributes {
			inventoryAttr := entities.ProductInventoryAttribute{
				ProductID:          productInventory.ProductID,
				ProductInventoryID: uint(inventoryID),
				ProductAttributeID: attr.ID,
			}
			if err := tx.Create(&inventoryAttr).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}
	return nil
}

func (p ProductRepository) UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) error {
	var inventory entities.ProductInventory
	if iErr := p.db.WithContext(c).First(&inventory, inventoryID).Error; iErr != nil {
		return iErr
	}

	if updateErr := p.db.WithContext(c).Model(&inventory).Update("quantity", quantity).Error; updateErr != nil {
		return updateErr
	}
	return nil
}

func (p ProductRepository) InsertFeature(c *gin.Context, productID int, req requests.CreateProductFeatureRequest) error {

	if err := p.db.Create(&entities.Feature{
		ProductID: uint(productID),
		Title:     req.Title,
		Value:     req.Value,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (p ProductRepository) DeleteFeature(c *gin.Context, productID int, featureID int) error {
	if err := p.db.WithContext(c).Where("product_id=? ", productID).Where("id = ?", featureID).Unscoped().Delete(&entities.Feature{}).Error; err != nil {
		return err
	}
	return nil
}

func (p ProductRepository) GetFeatureBy(c *gin.Context, productID int, featureID int) (entities.Feature, error) {
	var feature entities.Feature
	if err := p.db.WithContext(c).Where("id=?", featureID).Where("product_id=?", productID).First(&feature).Error; err != nil {
		return feature, err
	}
	return feature, nil
}

func (p ProductRepository) EditFeature(c *gin.Context, productID int, featureID int, req requests.UpdateProductFeatureRequest) error {
	if err := p.db.Where("id=?", featureID).Where("product_id=?", productID).Model(&entities.Feature{}).Update("title", req.Title).
		Update("value", req.Value).Error; err != nil {
		return err
	}
	return nil
}
