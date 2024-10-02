package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	CustomerID  uint
	ProductID   uint
	InventoryID uint

	Count  uint8 `gorm:"type:smallInt"`
	Status uint8

	//--
	ProductSku   string
	ProductTitle string
	ProductImage string
	ProductSlug  string

	OriginalPrice uint
	SalePrice     uint

	//----
	//Customer Customer `gorm:"references:ID"`
	//Product  Product  `gorm:"references:ID"`
}
