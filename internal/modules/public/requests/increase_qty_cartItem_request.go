package requests

type IncreaseCartItemQty struct {
	CartID      uint `form:"cart_id" binding:"required"`
	ProductID   uint `form:"product_id" binding:"required"`
	InventoryID uint `form:"inventory_id" binding:"required"`
}
