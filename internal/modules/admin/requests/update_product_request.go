package requests

type UpdateProductRequest struct {
	CategoryID    int    `form:"category_id" binding:"required"`
	BrandID       uint   `form:"brand_id" binding:"required"`
	Title         string `form:"title" binding:"required"`
	Slug          string `form:"slug" binding:"required"`
	Sku           string `form:"sku" binding:"required"`
	Status        string `form:"status"`
	OriginalPrice uint   `form:"original_price" binding:"required"`
	SalePrice     uint   `form:"sale_price" binding:"required"`
	Description   string `form:"description" binding:"required"`
}
