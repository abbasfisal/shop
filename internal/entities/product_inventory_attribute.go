package entities

import (
	"gorm.io/gorm"
)

type ProductInventoryAttribute struct {
	gorm.Model
	ProductID          uint
	ProductInventoryID uint
	ProductAttributeID uint

	//relations
	Product          Product          `gorm:"foreignKey:ProductID"`
	ProductInventory ProductInventory `gorm:"foreignKey:ProductInventoryID"`
	ProductAttribute ProductAttribute `gorm:"foreignKey:ProductAttributeID"`
}
