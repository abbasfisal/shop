package responses

import "shop/internal/entities"

type Attribute struct {
	ID         uint
	CategoryID uint
	Title      string
}
type Attributes struct {
	Data []Attribute
}

func ToAttribute(attr entities.Attribute) Attribute {
	return Attribute{
		ID:         attr.ID,
		CategoryID: attr.CategoryID,
		Title:      attr.Title,
	}
}

func ToAttributes(attr []entities.Attribute) Attributes {
	var response Attributes

	for _, item := range attr {
		response.Data = append(response.Data, ToAttribute(item))
	}

	return response
}
