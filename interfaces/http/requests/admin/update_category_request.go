package requests

type UpdateCategoryRequest struct {
	CategoryID uint   `form:"category_id"`
	Title      string `form:"title" binding:"required"`
	Priority   *uint  `form:"priority"`
	Slug       string `form:"slug" binding:"required"`
	Status     string `form:"status"`
	Image      string
}
