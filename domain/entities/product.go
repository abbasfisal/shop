package entities

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryID    uint
	BrandID       uint
	Title         string
	Slug          string `gorm:"unique"`
	Sku           string `gorm:"unique"`
	Status        bool
	OriginalPrice uint
	SalePrice     uint
	Description   string

	// flattened read model (JSONB) - replaces the old MongoDB products collection
	ReadModel datatypes.JSON `json:"-" gorm:"type:jsonb;default:'{}'"`

	//--------------relations
	///////////////////////////////////

	Category                   *Category                    `gorm:"foreignKey:CategoryID"`
	Brand                      *Brand                       `gorm:"foreignKey:BrandID"`
	ProductImages              []*ProductImages             `gorm:"foreignKey:ProductID"`
	ProductAttributes          []*ProductAttribute          `gorm:"foreignKye:ProductID"`
	ProductInventories         []*ProductInventory          `gorm:"foreignKey:ProductID"`
	ProductInventoryAttributes []*ProductInventoryAttribute `gorm:"foreignKey:ProductID"`
	Features                   []*Feature                   `gorm:"foreignKey:ProductID"`
	///Carts                      []Cart                      `gorm:"foreignKey:ProductID"`
}
