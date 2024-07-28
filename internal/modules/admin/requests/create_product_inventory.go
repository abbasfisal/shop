package requests

type CreateProductInventoryRequest struct {
	ProductAttributes []uint `form:"productAttributes" binding:"required"`
	Quantity          uint   `form:"quantity" binding:"required"`
}
