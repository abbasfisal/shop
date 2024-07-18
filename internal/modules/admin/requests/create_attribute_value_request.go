package requests

type CreateAttributeValueRequest struct {
	AttributeID uint   `form:"attribute_id" binding:"required"`
	Value       string `form:"value" binding:"required"`
}
