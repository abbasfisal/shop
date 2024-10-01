package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID      uint
	ProductID   uint
	InventoryID uint

	Count  uint8 `gorm:"type:smallInt"`
	Status uint8

	//--
	ProductTitle string
	ProductImage string
	ProductSlug  string

	OriginalPrice uint
	SalePrice     uint

	//----

	Product Product
}

//user1
//product1
//inventory1
//count2
//
