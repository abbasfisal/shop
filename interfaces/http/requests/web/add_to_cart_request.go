package requests

type AddToCartRequest struct {
	ProductID   string `form:"product_id" binding:"required"` //ProductID is a ObjectID in mongo document
	InventoryID uint   `form:"inventory_id"`
}
