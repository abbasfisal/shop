package entities

import (
	"time"

	"gorm.io/gorm"
)

// Variant status values (Laravel VariantStatus enum pattern).
const (
	VariantStatusActive   = "active"
	VariantStatusInactive = "inactive"
)

// ProductVariant is one sellable combination of a product (Laravel ProductVariant
// pattern): its own stock, reservation and optional per-variant prices.
// A NULL price field means "inherit the parent product's price" — the existing
// storefront inventory form only manages stock, so overrides are applied via
// code/seed/API while aggregates always resolve the effective price.
type ProductVariant struct {
	gorm.Model
	ProductID     uint
	Sku           *string    `gorm:"unique"`
	Price         *uint      // buy/list price; NULL = inherit products.original_price
	SalePrice     *uint      // NULL = inherit products.sale_price
	DiscountPrice *uint      // counts only when > 0 and < effective sale price
	Stock         uint       `gorm:"default:0"`
	ReservedStock uint       `gorm:"default:0"`
	Status        string     `gorm:"type:varchar(16);default:'active'"`
	ExpiresAt     *time.Time `gorm:"type:date"`

	Product                Product                  `gorm:"foreignKey:ProductID"`
	VariantAttributeValues []*VariantAttributeValue `gorm:"foreignKey:VariantID"`
}

// VariantAttributeValue links a variant to an attribute value
// (Laravel variant_attribute_values junction). product_id is denormalized so
// order-detail preloads can filter like the legacy join table did; a
// BeforeCreate hook fills it when only the variant is known.
type VariantAttributeValue struct {
	gorm.Model
	ProductID        uint
	VariantID        uint `gorm:"index"`
	AttributeValueID uint `gorm:"index"`

	Variant        ProductVariant  `gorm:"foreignKey:VariantID"`
	AttributeValue *AttributeValue `gorm:"foreignKey:AttributeValueID"`
	Product        Product         `gorm:"foreignKey:ProductID"`
}

func (v *VariantAttributeValue) BeforeCreate(tx *gorm.DB) error {
	if v.ProductID == 0 && v.VariantID > 0 {
		var pv ProductVariant
		if err := tx.Select("product_id").First(&pv, v.VariantID).Error; err == nil {
			v.ProductID = pv.ProductID
		}
	}
	return nil
}
