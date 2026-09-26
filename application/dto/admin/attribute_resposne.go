package responses

import "shop/domain/entities"

type Attribute struct {
	ID              uint
	Title           string
	Code            string
	SortOrder       int
	InputType       string
	InputTypeLabel  string
	AttributeValues *AttributeValues
}

type Attributes struct {
	Data []Attribute
}

func ToAttribute(attr *entities.Attribute) *Attribute {
	inputType := attr.InputType
	if inputType == "" {
		inputType = entities.AttributeInputText
	}
	return &Attribute{
		ID:              attr.ID,
		Title:           attr.Title,
		Code:            attr.Code,
		SortOrder:       attr.SortOrder,
		InputType:       inputType,
		InputTypeLabel:  entities.AttributeInputTypeLabel(inputType),
		AttributeValues: ToAttributeValues(attr.AttributeValues),
	}
}

func ToAttributes(attr []*entities.Attribute) *Attributes {
	if attr == nil {
		return &Attributes{}
	}

	var response Attributes
	for _, item := range attr {
		response.Data = append(response.Data, *ToAttribute(item))
	}

	return &response
}
