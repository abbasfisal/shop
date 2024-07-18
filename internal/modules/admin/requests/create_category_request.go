package requests

type CreateCategoryRequest struct {
	ParentID uint   `form:"parent_id"`
	Title    string `form:"title" binding:"required"`
	Slug     string `form:"slug" binding:"required"`
	Status   string `form:"status"`
	Image    string
}
