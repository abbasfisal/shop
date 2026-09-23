package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

// StoreProductInventory creates a product variant (stock row) and links the
// selected product attribute-values to it via variant_attribute_values.
// The HTML inventory form only posts productAttributes + quantity; prices are
// left NULL so the variant inherits the product price (see PricingService).
func (p *ProductRepository) StoreProductInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) (*entities.ProductVariant, error) {

	var variant entities.ProductVariant

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var productAttributes []entities.ProductAttribute

		//fetch product-attributes
		//len(req.ProductAttributes)<=0 means: stock-only variant without attributes
		if len(req.ProductAttributes) > 0 {
			if err := tx.WithContext(c).Where("id IN ? ", req.ProductAttributes).Find(&productAttributes).Error; err != nil {
				return err
			}
			//check len retrieved product-attribute
			if len(productAttributes) != len(req.ProductAttributes) {
				return gorm.ErrRecordNotFound
			}
		}

		variant = entities.ProductVariant{
			ProductID: uint(productID),
			Stock:     req.Quantity,
			Status:    entities.VariantStatusActive,
		}

		if iErr := tx.WithContext(c).Create(&variant).Error; iErr != nil {
			return iErr
		}

		//link selected attribute values to the new variant
		if len(req.ProductAttributes) > 0 {
			for _, attr := range productAttributes {
				vav := entities.VariantAttributeValue{
					ProductID:        uint(productID),
					VariantID:        variant.ID,
					AttributeValueID: attr.AttributeValueID,
				}
				if err := tx.Create(&vav).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if txErr != nil {
		fmt.Println("---- create product variant err: ", txErr)
		return nil, txErr
	}

	_ = SyncReadModel(c, p.db, uint(productID))

	return &variant, nil
}

// DeleteInventoryAttribute removes one variant <-> attribute_value link.
// Returns the affected product id for pricing refresh.
func (p *ProductRepository) DeleteInventoryAttribute(c *gin.Context, variantAttributeValueID int) (uint, error) {
	var vav entities.VariantAttributeValue
	if err := p.db.WithContext(c).First(&vav, variantAttributeValueID).Error; err != nil {
		return 0, err
	}

	//hard delete so the UNIQUE(variant_id, attribute_value_id) slot can be reused
	if piaErr := p.db.WithContext(c).Unscoped().Delete(&vav).Error; piaErr != nil {
		return 0, piaErr
	}

	_ = SyncReadModel(c, p.db, vav.ProductID)

	return vav.ProductID, nil
}

// DeleteInventory soft-deletes a variant and hard-deletes its attribute links.
// Returns the affected product id for pricing refresh.
func (p *ProductRepository) DeleteInventory(c *gin.Context, inventoryID int) (uint, error) {
	var productID uint

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var variant entities.ProductVariant

		if iErr := tx.WithContext(c).First(&variant, inventoryID).Error; iErr != nil {
			return iErr
		}

		productID = variant.ProductID

		//hard-delete attribute links (frees the UNIQUE slot)
		if deleteErr := tx.Where("variant_id = ?", variant.ID).
			Unscoped().
			Delete(&entities.VariantAttributeValue{}).Error; deleteErr != nil {
			return deleteErr
		}

		if iDelete := tx.WithContext(c).Delete(&variant).Error; iDelete != nil {
			return iDelete
		}
		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	_ = SyncReadModel(c, p.db, productID)

	return productID, nil
}

// AppendAttributesToInventory links additional attribute values to a variant.
// Returns the affected product id for pricing refresh.
func (p *ProductRepository) AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) (uint, error) {
	var variant entities.ProductVariant

	if err := p.db.WithContext(c).First(&variant, inventoryID).Error; err != nil {
		return 0, err
	}

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var productAttributes []entities.ProductAttribute

		if err := tx.WithContext(c).Where("id IN ? ", attributes).Find(&productAttributes).Error; err != nil {
			return err
		}
		if len(productAttributes) != len(attributes) {
			return gorm.ErrRecordNotFound
		}

		//existing links (UNIQUE constraint would reject duplicates)
		var existing []uint
		if err := tx.Model(&entities.VariantAttributeValue{}).
			Where("variant_id = ?", variant.ID).
			Pluck("attribute_value_id", &existing).
			Error; err != nil {
			return err
		}
		exists := make(map[uint]struct{}, len(existing))
		for _, id := range existing {
			exists[id] = struct{}{}
		}

		for _, attr := range productAttributes {
			if _, ok := exists[attr.AttributeValueID]; ok {
				continue
			}
			vav := entities.VariantAttributeValue{
				ProductID:        variant.ProductID,
				VariantID:        variant.ID,
				AttributeValueID: attr.AttributeValueID,
			}
			if err := tx.Create(&vav).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	_ = SyncReadModel(c, p.db, variant.ProductID)

	return variant.ProductID, nil
}

// UpdateInventoryQuantity updates the variant stock.
// Returns the affected product id for pricing refresh.
func (p *ProductRepository) UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) (uint, error) {
	var variant entities.ProductVariant
	if iErr := p.db.WithContext(c).First(&variant, inventoryID).Error; iErr != nil {
		return 0, iErr
	}

	if updateErr := p.db.WithContext(c).Model(&variant).Update("stock", quantity).Error; updateErr != nil {
		return 0, updateErr
	}

	_ = SyncReadModel(c, p.db, variant.ProductID)

	return variant.ProductID, nil
}
