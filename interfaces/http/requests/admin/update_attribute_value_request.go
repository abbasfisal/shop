package requests

type UpdateAttributeValueRequest struct {
	AttributeID    uint `form:"attribute_id" binding:"required,gt=0"`
	AttributeTitle string
	Value          string `form:"value" binding:"required"`
	// hex color (#rrggbb) — only meaningful when the parent attribute is color
	ColorHex string `form:"color_hex"`
}
