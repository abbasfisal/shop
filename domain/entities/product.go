package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product status values (Laravel ProductStatus enum pattern).
const (
	ProductStatusDraft     = "draft"
	ProductStatusPublished = "published"
	ProductStatusArchived  = "archived"
)

// ProductStatusLabel is the Persian label of a product status (admin UI).
func ProductStatusLabel(status string) string {
	switch status {
	case ProductStatusPublished:
		return "منتشر شده"
	case ProductStatusArchived:
		return "بایگانی"
	default:
		return "پیش‌نویس"
	}
}

// ValidProductStatus reports whether status is one of the enum values.
func ValidProductStatus(status string) bool {
	switch status {
	case ProductStatusDraft, ProductStatusPublished, ProductStatusArchived:
		return true
	}
	return false
}

type Product struct {
	gorm.Model
	CategoryID    uint
	BrandID       uint
	Title         string
	Slug          string `gorm:"unique"`
	Sku           string `gorm:"unique"`
	Status        string `gorm:"type:varchar(16);default:'draft'"`
	OriginalPrice uint
	SalePrice     uint
	Description   string
	ExpiresAt     *time.Time `gorm:"type:date"`

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
