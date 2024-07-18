package requests

type CreateAttributeRequest struct {
	CategoryID uint   `form:"category_id" binding:"required"`
	Title      string `form:"title" binding:"required"`
}
