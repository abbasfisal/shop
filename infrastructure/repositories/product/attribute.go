package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"strconv"
	"strings"
)

func (p *ProductRepository) GetRootAttributes(c *gin.Context, productID int) ([]*entities.Attribute, error) {
	var category entities.Category
	var attributes []*entities.Attribute

	var product entities.Product
	pErr := p.db.WithContext(c).Where("id = ? ", productID).First(&product).Error

	if pErr != nil {
		return nil, pErr
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
		return nil, err
	}

	aErr := p.db.WithContext(c).Preload("AttributeValues").Find(&attributes).Error

	return attributes, aErr
}

func (p *ProductRepository) StoreAttributeValues(ctx *gin.Context, productID int, attValues []string) error {
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
	_ = SyncReadModel(ctx, p.db, uint(productID))
	return nil
}

func (p *ProductRepository) GetProductAndAttributes(ctx *gin.Context, productID int) (map[string]interface{}, error) {
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
	aerr := p.db.WithContext(ctx).Where("id = ?", productID).First(&product).Error

	if aerr != nil {
		return map[string]interface{}{}, aerr
	}

	var inventories []InventoryWithAttributes

	result := make(map[string]interface{})

	serr := p.db.WithContext(ctx).
		Table("product_variants").
		Select("product_variants.id AS inventory_id, product_variants.stock AS quantity, attributes.id AS attribute_id, attributes.title AS attribute_title, attribute_values.id AS attribute_value_id, attribute_values.value AS attribute_value_title, variant_attribute_values.id AS product_inventory_attribute_id").
		Joins("LEFT JOIN variant_attribute_values ON product_variants.id = variant_attribute_values.variant_id AND variant_attribute_values.deleted_at IS NULL").
		Joins("LEFT JOIN attribute_values ON variant_attribute_values.attribute_value_id = attribute_values.id AND attribute_values.deleted_at IS NULL").
		Joins("LEFT JOIN attributes ON attribute_values.attribute_id = attributes.id AND attributes.deleted_at IS NULL").
		Where("product_variants.product_id = ? and product_variants.deleted_at IS NULL", productID).
		Scan(&inventories).
		Error

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

// DeleteProductAttribute soft-deletes one product_attributes row (legacy
// "add attributes" page). Variant links live in variant_attribute_values, so
// the pricing aggregates are not affected — only the read model is rebuilt.
func (p *ProductRepository) DeleteProductAttribute(c *gin.Context, productAttributeID int) (uint, error) {
	var pa entities.ProductAttribute
	if err := p.db.WithContext(c).First(&pa, productAttributeID).Error; err != nil {
		return 0, err
	}
	if err := p.db.WithContext(c).Delete(&pa).Error; err != nil {
		return 0, err
	}

	_ = SyncReadModel(c, p.db, pa.ProductID)
	return pa.ProductID, nil
}
