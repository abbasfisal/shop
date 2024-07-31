package requests

type CreateBrandRequest struct {
	Title string `form:"title" binding:"required"`
	Slug  string `form:"slug" binding:"required"`
	Image string
}
