package requests

type UpdateBrandRequest struct {
	Title string `form:"title"`
	Slug  string `form:"slug"`
	Image string
}
