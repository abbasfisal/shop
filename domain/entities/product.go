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

	// aggregate cache refreshed by PricingService (Laravel pattern)
	ProductType    string         `gorm:"column:product_type;type:varchar(16);default:'simple'"`
	MinPrice       uint           `gorm:"default:0"`
	MaxPrice       uint           `gorm:"default:0"`
	TotalStock     uint           `gorm:"default:0"`
	TotalReserved  uint           `gorm:"default:0"`
	AvailableStock uint           `gorm:"default:0"`
	InStock        bool           `gorm:"default:false"`
	VariantsCount  int            `gorm:"default:0"`
	AttributesJSON datatypes.JSON `gorm:"column:attributes_json;type:jsonb;default:'{}'"`

	//--------------relations
	///////////////////////////////////

	Category               *Category                `gorm:"foreignKey:CategoryID"`
	Brand                  *Brand                   `gorm:"foreignKey:BrandID"`
	ProductImages          []*ProductImages         `gorm:"foreignKey:ProductID"`
	ProductAttributes      []*ProductAttribute      `gorm:"foreignKey:ProductID"`
	ProductVariants        []*ProductVariant        `gorm:"foreignKey:ProductID"`
	VariantAttributeValues []*VariantAttributeValue `gorm:"foreignKey:ProductID"`
	Features               []*Feature               `gorm:"foreignKey:ProductID"`
}
