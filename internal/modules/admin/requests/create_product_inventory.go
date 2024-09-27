package requests

type CreateProductInventoryRequest struct {
	ProductAttributes []uint `form:"productAttributes" `
	Quantity          uint   `form:"quantity" binding:"required"`
}
