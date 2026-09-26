package requests

import "shop/domain/entities"

type CreateAttributeRequest struct {
	Title     string `form:"title" binding:"required"`
	InputType string `form:"input_type"`
}

// NormalizedInputType defaults an empty/unknown type to plain text.
func (r *CreateAttributeRequest) NormalizedInputType() string {
	if entities.ValidAttributeInputType(r.InputType) {
		return r.InputType
	}
	return entities.AttributeInputText
}
