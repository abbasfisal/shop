package entities

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CustomerID    uint
	CartID        uint
	ProductID     uint
	InventoryID   uint
	Quantity      uint8
	OriginalPrice uint
	SalePrice     uint

	ProductSku   string
	ProductTitle string
	ProductImage string
	ProductSlug  string

	//---
}
