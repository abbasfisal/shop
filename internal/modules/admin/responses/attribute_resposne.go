package responses

import "shop/internal/entities"

type Attribute struct {
	ID              uint
	Title           string
	AttributeValues *AttributeValues
}

type Attributes struct {
	Data []Attribute
}

func ToAttribute(attr *entities.Attribute) *Attribute {
	return &Attribute{
		ID:              attr.ID,
		Title:           attr.Title,
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
