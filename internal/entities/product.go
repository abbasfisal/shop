package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryID uint
	BrandID    uint
	Title      string
	Slug       string `gorm:"unique"`
	Sku        string `gorm:"unique"`
	Status     bool
	//Quantity      uint `gorm:"not null"`
	OriginalPrice uint
	SalePrice     uint
	Description   string
	Category      Category        `gorm:"foreignKey:CategoryID"`
	Brand         Brand           `gorm:"foreignKey:BrandID"`
	ProductImages []ProductImages `gorm:"foreignKey:ProductID"`

	ProductAttributes  []ProductAttribute `gorm:"foreignKye:ProductID"`
	ProductInventories []ProductInventory `gorm:"foreignKey:ProductID"`
}
