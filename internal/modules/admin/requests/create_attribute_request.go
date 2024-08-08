package requests

type CreateAttributeRequest struct {
	Title string `form:"title" binding:"required"`
}
