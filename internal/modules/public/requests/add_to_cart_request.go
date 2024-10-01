package requests

type AddToCartRequest struct {
	ProductID   string `form:"product_id" binding:"required"`
	InventoryID *uint  `form:"inventory_id"`
}
