package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/requests"
)

func (p *ProductRepository) StoreProductInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) (*entities.ProductInventory, error) {

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
			Quantity:  req.Quantity,
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
		return nil, txErr
	}

	_ = SyncMongo(c, p.db, uint(productID))

	return &inventory, nil
}

func (p *ProductRepository) DeleteInventoryAttribute(c *gin.Context, productInventoryAttributeID int) error {

	//find
	var productInventoryAttribute entities.ProductInventoryAttribute
	if err := p.db.First(&productInventoryAttribute, productInventoryAttributeID).Error; err != nil {
		return err
	}

	//delete from product_inventory_attributes table
	if piaErr := p.db.WithContext(c).Unscoped().Delete(&productInventoryAttribute).Error; piaErr != nil {
		return piaErr
	}

	_ = SyncMongo(c, p.db, productInventoryAttribute.ProductID)

	return nil
}

func (p *ProductRepository) DeleteInventory(c *gin.Context, inventoryID int) error {

	var productID uint

	txErr := p.db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		var inventory entities.ProductInventory

		//find inventory
		if iErr := p.db.WithContext(c).First(&inventory, inventoryID).Error; iErr != nil {
			return iErr
		}

		productID = inventory.ProductID

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

	_ = SyncMongo(c, p.db, productID)

	return nil
}

func (p *ProductRepository) AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) error {

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

	_ = SyncMongo(c, p.db, productInventory.ProductID)

	return nil
}

func (p *ProductRepository) UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) error {
	var inventory entities.ProductInventory
	if iErr := p.db.WithContext(c).First(&inventory, inventoryID).Error; iErr != nil {
		return iErr
	}

	if updateErr := p.db.WithContext(c).Model(&inventory).Update("quantity", quantity).Error; updateErr != nil {
		return updateErr
	}

	_ = SyncMongo(c, p.db, inventory.ProductID)

	return nil
}
