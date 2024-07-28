package requests

type CreateAttributeValueRequest struct {
	AttributeID uint   `form:"attribute_id" binding:"required,gt=0"`
	Value       string `form:"value" binding:"required"`
}
