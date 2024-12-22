package responses

import (
	"shop/internal/entities"
)

type ProductInventoryAttribute struct {
	ID                 uint
	ProductID          uint
	ProductInventoryID uint
	ProductAttributeID uint
	ProductAttribute   *ProductAttribute
}

type ProductInventoryAttributes struct {
	Data []ProductInventoryAttribute
}

func ToProductInventoryAttribute(pi *entities.ProductInventoryAttribute) *ProductInventoryAttribute {
	return &ProductInventoryAttribute{
		ID:        pi.ID,
		ProductID: pi.ProductID,

		ProductInventoryID: pi.ProductInventoryID,
		ProductAttributeID: pi.ProductAttributeID,

		//
		ProductAttribute: ToProductAttribute(pi.ProductAttribute),
	}
}

func ToProductInventoryAttributes(pis []*entities.ProductInventoryAttribute) *ProductInventoryAttributes {
	var pResponse ProductInventoryAttributes
	for _, p := range pis {
		pResponse.Data = append(pResponse.Data, *ToProductInventoryAttribute(p))
	}
	return &pResponse
}
