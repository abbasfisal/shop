package responses

import "shop/internal/entities"

type ProductAttribute struct {
	ProductID uint

	AttributeID    uint
	AttributeTitle string

	AttributeValueID    uint
	AttributeValueTitle string
}
type ProductAttributes struct {
	Data []ProductAttribute
}

func ToProductAttribute(productAttribute entities.ProductAttribute) ProductAttribute {
	return ProductAttribute{
		ProductID:           productAttribute.ProductID,
		AttributeID:         productAttribute.AttributeID,
		AttributeTitle:      productAttribute.AttributeTitle,
		AttributeValueID:    productAttribute.AttributeValueID,
		AttributeValueTitle: productAttribute.AttributeValueTitle,
	}
}

func ToProductAttributes(productAttributes []entities.ProductAttribute) ProductAttributes {
	var response ProductAttributes

	for _, productAtt := range productAttributes {
		response.Data = append(response.Data, ToProductAttribute(productAtt))
	}

	return response
}
