package requests

type CreateProductFeatureRequest struct {
	Title string `form:"title" binding:"required"`
	Value string `form:"value" binding:"required"`
}
