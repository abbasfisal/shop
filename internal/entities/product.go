package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryID uint
	Title      string
	Slug       string `gorm:"unique"`
	Sku        string `gorm:"unique"`
	Status     bool
	//Quantity      uint `gorm:"not null"`
	OriginalPrice uint
	SalePrice     uint
	Description   string
	Category      Category        `gorm:"foreignKey:CategoryID"`
	ProductImage  []ProductImages `gorm:"foreignKey:ProductID"`

	ProductAttribute []ProductAttribute `gorm:"foreignKye:ProductID"`
}
