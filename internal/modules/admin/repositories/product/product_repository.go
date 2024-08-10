package product

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
	"strconv"
	"strings"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db: db}
}

func (p ProductRepository) GetAll(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	err := p.db.WithContext(ctx).Preload("Category").Preload("Brand").Find(&products).Error
	return products, err
}

func (p ProductRepository) FindBy(ctx context.Context, columnName string, value any) (entities.Product, error) {
	var product entities.Product
	condition := fmt.Sprintf("%s=?", columnName)
	err := p.db.Preload("Category").Preload("Brand").Preload("ProductAttributes").Preload("ProductInventories").Preload("ProductImages").First(&product, condition, value).Error
	return product, err
}

func (p ProductRepository) FindByID(ctx context.Context, ID int) (entities.Product, error) {
	var product entities.Product
	err := p.db.WithContext(ctx).Preload("ProductAttributes").First(&product, ID).Error
	return product, err
}

func (p ProductRepository) Store(ctx context.Context, product entities.Product) (entities.Product, error) {

	err := p.db.WithContext(ctx).Create(&product).Error
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
		Joins("LEFT JOIN product_inventory_attributes ON product_inventories.id = product_inventory_attributes.product_inventory_id").
		Joins("LEFT JOIN product_attributes ON product_inventory_attributes.product_attribute_id = product_attributes.id").
		Joins("LEFT JOIN attributes ON product_attributes.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_values ON product_attributes.attribute_value_id = attribute_values.id").
		Where("product_inventories.product_id = ?", productID).
		Scan(&inventories).Error

	if serr != nil {
		return map[string]interface{}{}, serr
	}

	inventoryMap := make(map[uint]map[string]interface{})
	for _, inventory := range inventories {
		if _, exists := inventoryMap[inventory.InventoryID]; !exists {
			inventoryMap[inventory.InventoryID] = map[string]interface{}{
				"quantity":         inventory.Quantity,
				"delete_inventory": fmt.Sprintf("/admins/inventories/%d/delete", inventory.InventoryID), //remove record from product_inventories table
				"edit_inventory":   fmt.Sprintf("/admins/inventories/%d/edit", inventory.InventoryID),   //edit quantity of a product inventory (product_inventories)
				"inventory_id":     inventory.InventoryID,
				"attributes":       []map[string]interface{}{},
			}
		}

		attributes := inventoryMap[inventory.InventoryID]["attributes"].([]map[string]interface{})
		attributes = append(attributes, map[string]interface{}{
			"attribute_id":                   inventory.AttributeID,
			"attribute_title":                inventory.AttributeTitle,
			"attribute_value_id":             inventory.AttributeValueID,
			"attribute_value_title":          inventory.AttributeValueTitle,
			"product_inventory_attribute_id": inventory.ProductInventoryAttributeID,
			"delete_attribute":               fmt.Sprintf("/admins/inventories/%d/attributes/delete", inventory.ProductInventoryAttributeID), //remove from product_inventory_attributes
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
		if err := p.db.WithContext(c).Where("id IN ? ", req.ProductAttributes).Find(&productAttributes).Error; err != nil {
			return err
		}
		//check len retrieved product-attribute
		if len(productAttributes) != len(req.ProductAttributes) {
			return gorm.ErrRecordNotFound
		}

		inventory = entities.ProductInventory{
			ProductID: uint(productID),
			Quantity:  uint(req.Quantity),
		}

		//store inventory
		if iErr := tx.Create(&inventory).Error; iErr != nil {
			return iErr
		}

		//store product-attribute in product-inventory-attribute table
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
