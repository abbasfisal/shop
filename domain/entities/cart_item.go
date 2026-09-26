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

	// Attributes are the labels of the selected variant («رنگ: قرمز»,
	// «سایز: M») resolved from variant_attribute_values when the cart is
	// loaded — not a stored column.
	Attributes []CartItemAttribute `gorm:"-"`

	//---
}

// CartItemAttribute is one attribute of the variant in the cart line.
type CartItemAttribute struct {
	Title string
	Value string
}
