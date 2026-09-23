package requests

type AddToCartRequest struct {
	ProductID   string `form:"product_id" binding:"required"` // ProductID is the numeric product id (string form value)
	InventoryID uint   `form:"inventory_id"`
}
