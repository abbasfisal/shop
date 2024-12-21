package responses

import "shop/internal/entities"

type AttributeValue struct {
	ID             uint
	AttributeID    uint
	AttributeTitle string
	Title          string
}

type AttributeValues struct {
	Data []AttributeValue
}

func ToAttributeValue(attr *entities.AttributeValue) AttributeValue {
	return AttributeValue{
		ID:             attr.ID,
		AttributeID:    attr.AttributeID,
		AttributeTitle: attr.AttributeTitle,
		Title:          attr.Value,
	}
}

func ToAttributeValues(attrValue []entities.AttributeValue) *AttributeValues {
	var response AttributeValues

	for _, item := range attrValue {
		response.Data = append(response.Data, ToAttributeValue(&item))
	}

	return &response // Return a value (not pointer)
}
