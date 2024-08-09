package entities

import (
	"gorm.io/gorm"
)

type ProductInventory struct {
	gorm.Model
	ProductID uint
	Quantity  uint
}
