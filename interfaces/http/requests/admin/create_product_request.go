package requests

type CreateProductRequest struct {
	CategoryID    int    `form:"category_id" binding:"required"`
	BrandID       uint   `form:"brand_id" binding:"required"`
	Title         string `form:"title" binding:"required"`
	Slug          string `form:"slug" binding:"required"`
	Sku           string `form:"sku" binding:"required"`
	Status        string `form:"status"`
	ExpiresAt     string `form:"expires_at"`
	OriginalPrice uint   `form:"original_price" binding:"required"`
	SalePrice     uint   `form:"sale_price" binding:"required"`
	Description   string `form:"description" binding:"required"`

	// simple | variable (Laravel ProductType enum)
	ProductType string `form:"product_type"`

	// product_type = simple → single variant built from the base price rows
	SimpleDiscountPrice *uint  `form:"simple_discount_price"`
	SimpleStock         string `form:"simple_stock"`
	SimpleStatus        string `form:"simple_status"`
	SimpleExpiresAt     string `form:"simple_expires_at"`

	ProductImage []string

	// parsed from variants[i][...] keys (empty for simple products)
	Variants []VariantRow
}
